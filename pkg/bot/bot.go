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
	"github.com/RomanMalashenkov/tg_bot/pkg/filehandler"

	"github.com/sunshineplan/imgconv"
	tele "gopkg.in/telebot.v3"
)

var (
	//fileURL    string
	userFileID string
	userFormat imgconv.Format
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

	log.Print("Бот готов к работе")

	// ответ на команду /start
	b.Handle("/start", func(c tele.Context) error {
		log.Printf("[LOG]: User: %s | Controller: /start ", c.Message().Sender.Username)
		_, err := b.Send(c.Chat(), "Привет, я бот-конвертер\nПришлите мне файл")
		return err
	})

	// ответ на команду /help
	b.Handle("/help", func(c tele.Context) error {
		log.Printf("[LOG]: User: %s | Controller: /help ", c.Message().Sender.Username)
		_, err := b.Send(c.Chat(), "Я умею работать со следующими форматами:\njpg, jpeg, png, gif, tif, tiff, bmp, pdf")
		return err
	})

	//var format imgconv.Format

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
			InlineKeyboard: convertToInlineButtons(btns),
		})
		if err != nil {
			return err
		}

		// Сохраняем информацию о файле для последующей обработки
		//userFileID = c.Message().Document.FileID
		userFileID = doc.FileID
		/*
			// Получаем URL файла по его FileID
			fileURL, err = api.GetFileURL(b, doc.FileID)
			if err != nil {
				log.Printf("Failed to get file URL: %s", err)
				return err
			}

			// Отправляем запрос на конвертацию изображения
			err = filehandler.ConvertAndSendImage(fileURL, c, b, format)
			if err != nil {
				log.Printf("Failed to convert and send image: %s", err)
				return err
			}
		*/
		return nil
	})
	////////////////////////////////////////////////////////////////////////////////////
	// обработка выбора формата конвертации
	b.Handle(tele.OnCallback, func(c tele.Context) error {
		data := c.Callback().Data // появляется префикс ♀

		log.Printf("Получены данные из Callback: %v", data)

		if userFileID != "" {
			fileURL, err := api.GetFileURL(b, userFileID)
			if err != nil {
				log.Printf("Failed to get file URL: %s", err)
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

			log.Println("Пользователь выбрал формат конвертации:", selectedFormat)
			/////////////////////////////////////
			// Если формат выбран, обрабатываем его дальше
			/*if selectedFormat != 0 {
			// Определите ID файла и URL, чтобы затем выполнить конвертацию и отправку файла
			fileID := getFileIDFromContext(c)
			fileURL, err := api.GetFileURL(b, fileID)
			if err != nil {
				log.Printf("Failed to get file URL: %s", err)
				return err
			}
			*/
			// Выполните конвертацию и отправку файла
			err = filehandler.ConvertAndSendImage(fileURL, c, b, selectedFormat)
			if err != nil {
				log.Printf("Failed to convert and send file: %s", err)
				return err
			}
		} else {
			log.Println("Сообщение не содержит документ")

		}

		/*
			switch data {
			case "convert_to_jpeg":
				userFormat = imgconv.JPEG
			case "convert_to_png":
				userFormat = imgconv.PNG
			case "convert_to_gif":
				userFormat = imgconv.GIF
			case "convert_to_tiff":
				userFormat = imgconv.TIFF
			case "convert_to_bmp":
				userFormat = imgconv.BMP
			case "convert_to_pdf":
				userFormat = imgconv.PDF
			}

			log.Printf("Пользователь выбрал формат конвертации %v", userFormat)

			// Если есть информация о файле и выбранном формате - конвертируем и отправляем файл
			if userFileID != "" && userFormat != 0 {
				fileURL, err := api.GetFileURL(b, userFileID)
				if err != nil {
					log.Printf("Failed to get file URL: %s", err)
					return err
				} else {
					log.Println("отпр ф через апи тг") ////////////////
				}

				err = filehandler.ConvertAndSendImage(fileURL, c, b, userFormat)
				if err != nil {
					log.Printf("Failed to convert and send image: %s", err)
					return err
				} else {
					log.Println("отправ сконв ф") ////////////////
				}

				// Очищаем сохраненные данные
				userFileID = ""
				userFormat = 0
				log.Println("очи сохр") ////////////////
			}
		*/
		return nil
	})

	b.Start()
}

func convertToInlineButtons(btns [][]tele.Btn) [][]tele.InlineButton {
	var inlineBtns [][]tele.InlineButton
	for _, btnRow := range btns {
		var inlineRow []tele.InlineButton
		for _, btn := range btnRow {
			inlineRow = append(inlineRow, tele.InlineButton{
				Unique: btn.Data,
				Text:   btn.Text,
			})
			log.Printf("Уникальный ID кнопки: %s, Текст кнопки: %s\n", btn.Data, btn.Text)
			log.Printf("Тип данных btn.Data: %T", btn.Data)
		}
		inlineBtns = append(inlineBtns, inlineRow)
	}
	return inlineBtns
}

// Функция получения ID файла из контекста
func getFileIDFromContext(c tele.Context) string {
	switch c.Message().Document {
	case nil:
		log.Println("Сообщение не содержит документ")
		return ""
	default:
		return c.Message().Document.FileID
	}
}
