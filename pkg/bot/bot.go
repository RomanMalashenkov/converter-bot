package bot

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/RomanMalashenkov/tg_bot/pkg/api"
	"github.com/RomanMalashenkov/tg_bot/pkg/config"
	"github.com/RomanMalashenkov/tg_bot/pkg/converter"
	"github.com/RomanMalashenkov/tg_bot/pkg/queue"

	"github.com/sunshineplan/imgconv"
	tele "gopkg.in/telebot.v3"
)

var (
	userFileID string
	taskQueue  *queue.Queue //для очереди
)

func StartBot() {
	botConf, confErr := config.GetConfig()
	if confErr != nil {
		log.Fatal("No config")
	}

	b, err := tele.NewBot(tele.Settings{
		Token:  botConf.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Print("[LOG]: Бот готов к работе")

	// ответ на команду /start
	b.Handle("/start", func(c tele.Context) error {
		log.Printf("[LOG]: User: %s | Controller: /start ", c.Message().Sender.Username)
		_, err := b.Send(c.Chat(), `Привет! Я ConvBot!
Со мной вы можете конвертировать файлы одного формата в другой.
На данный момент я конвертирую только изображения.

Отправьте мне изображение в виде файла (документом).

Для получения дополнительной информации нажмите на каманду /help
		`)
		return err
	})

	// ответ на команду /help
	b.Handle("/help", func(c tele.Context) error {
		log.Printf("[LOG]: User: %s | Controller: /help ", c.Message().Sender.Username)
		_, err := b.Send(c.Chat(), "Поддерживающиеся форматы:\n\njpg, jpeg, png, gif, tif, tiff, bmp, pdf")
		return err
	})

	// инициализируем очередь
	taskQueue = queue.NewQueue()

	// обработка изображения, отправленного пользователем
	b.Handle(tele.OnDocument, func(c tele.Context) error {
		log.Printf("[LOG]: User: %s | Controller: OnDocument ", c.Message().Sender.Username)

		// получение инфы о файле
		doc := c.Message().Document
		//расширение файла
		docExt := filepath.Ext(doc.FileName)

		// проверка поддержки расширения для конвертации
		supported := converter.IsSupported(docExt)
		if !supported {
			msg := fmt.Sprintf("Извините, формат файла %s не поддерживается для конвертации", docExt)
			_, err := b.Send(c.Chat(), msg)
			return err
		}

		// Добавляем задачу в очередь
		taskID := doc.FileID // Используем ID файла в качестве задачи
		err := taskQueue.AddTaskToQueue(taskID)
		if err != nil {
			log.Printf("Не удалось добавить задачу в очередь: %s", err)
			return err
		} else {
			log.Printf("Задача успешно добавлена в очередь. ID задачи: %s", taskID)
		}

		// Проверяем содержимое очереди после добавления задачи
		task, err := taskQueue.GetTaskFromQueue()
		if err != nil {
			log.Printf("Не удалось получить задачу из очереди: %s", err)
		} else {
			log.Printf("Получена задача из очереди: %s", task)
		}
		// Обновляем userFileID для текущего файла
		userFileID = taskID
		// Создание кнопок для выбора формата конвертации
		btns := [][]tele.Btn{
			{
				tele.Btn{Text: "jpeg", Data: "convert_to_jpeg"},
				tele.Btn{Text: "png", Data: "convert_to_png"},
				tele.Btn{Text: "gif", Data: "convert_to_gif"},
				tele.Btn{Text: "tiff", Data: "convert_to_tiff"},
				tele.Btn{Text: "bmp", Data: "convert_to_bmp"},
				tele.Btn{Text: "pdf", Data: "convert_to_pdf"},
			},
		}

		// Отправка сообщения с кнопками выбора формата конвертации
		_, err = b.Send(c.Chat(), "Выберите формат для конвертации:", &tele.ReplyMarkup{
			InlineKeyboard: converter.ConvertToInlineButtons(btns),
		})
		if err != nil {
			return err
		}

		// Сохраняем информацию о файле для последующей обработки
		userFileID = doc.FileID

		return nil
	})

	// обработка выбора формата конвертации (когда польз нажал на кнопку)
	b.Handle(tele.OnCallback, func(c tele.Context) error {
		data := c.Callback().Data // появляется префикс ♀

		//log.Printf("Получены данные из Callback: %v", data)
		log.Printf("[LOG]: User: %s | Controller: OnCalback ", c.Message().Sender.Username)

		if userFileID != "" {
			fileURL, err := api.GetFileURL(b, userFileID)
			if err != nil {
				log.Printf("Не удалось получить URL-адрес файла: %s", err)
				return err
			}

			// Создадим карту для соответствия символов и форматов
			formats := map[string]imgconv.Format{
				"convert_to_jpeg": imgconv.JPEG,
				"convert_to_png":  imgconv.PNG,
				"convert_to_gif":  imgconv.GIF,
				"convert_to_tiff": imgconv.TIFF,
				"convert_to_bmp":  imgconv.BMP,
				"convert_to_pdf":  imgconv.PDF,
			}

			// Проверяем соответствие данных форматам, игнорируя символ ♀
			var selectedFormat imgconv.Format
			for key, val := range formats {
				if strings.HasSuffix(data, key) {
					selectedFormat = val
					break
				}
			}

			log.Println("Пользователь выбрал формат для конвертации:", selectedFormat)

			// Выполните конвертацию и отправку файла
			err = converter.ConvertAndSendImage(fileURL, c, b, selectedFormat)
			if err != nil {
				log.Printf("Не удалось преобразовать и отправить файл: %s", err)
				return err
			}
		} else {
			log.Println("Сообщение не содержит документ")
		}

		return nil
	})

	b.Start()
}
