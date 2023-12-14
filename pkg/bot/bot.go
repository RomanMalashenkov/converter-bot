package bot

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/RomanMalashenkov/tg_bot/pkg/config"
	"github.com/RomanMalashenkov/tg_bot/pkg/converter"
	"github.com/sunshineplan/imgconv"
	tele "gopkg.in/telebot.v3"
)

func getFileURL(bot *tele.Bot, fileID string) (string, error) {
	fileInfo, err := bot.FileByID(fileID)
	if err != nil {
		return "", err
	}

	// Получаем информацию о файле
	filePath := fileInfo.FilePath

	// Формируем URL для скачивания файла
	// Для получения прямой ссылки используем URL Telegram Bot API
	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, filePath)

	return fileURL, nil
}

func StartBot() {
	botConf, confErr := config.GetConfig()
	if confErr != nil {
		log.Fatal("No config")
	}

	//filehandler.CheckFolder(botConf.Store)
	//filehandler.CheckFolder(botConf.Store + "webm")
	//filehandler.CheckFolder(botConf.Store + "mp4")

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
		_, err := b.Send(c.Chat(), "Привет, я бот-конвертер")
		return err
	})

	// ответ на команду /help
	b.Handle("/help", func(c tele.Context) error {
		log.Printf("[LOG]: User: %s | Controller: /help ", c.Message().Sender.Username)
		_, err := b.Send(c.Chat(), "Пришлите мне webm, а я вам - mp4") ////////////////////////////изменить ссобщ
		return err
	})

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

		// Получаем URL файла по его FileID
		fileURL, err := getFileURL(b, doc.FileID)
		if err != nil {
			log.Printf("Failed to get file URL: %s", err)
			return err
		}
		// Отправляем запрос на конвертацию изображения
		err = convertAndSendImage(fileURL, c, b)
		if err != nil {
			log.Printf("Failed to convert and send image: %s", err)
			return err
		}
		/*
			// Создание клавиатуры для выбора формата конвертации
			btns := [][]tele.Btn{
				{
					tele.Btn{
						Text: fmt.Sprintf("Конвертировать в .jpg"),
						Data: "convert_to_jpg",
					},
					tele.Btn{
						Text: fmt.Sprintf("Конвертировать в .png"),
						Data: "convert_to_png",
					},
					// Добавьте другие форматы, если необходимо
				},
			}

				// Отправка сообщения с кнопками выбора формата конвертации
				_, err := b.Send(c.Chat(), "Выберите формат для конвертации:", &tele.ReplyMarkup{
					InlineKeyboard: btns,
				})
				if err != nil {
					return err
				}
		*/
		return nil
	})

	// обработка выбора формата конвертации

	b.Start()
}

func convertAndSendImage(fileURL string, c tele.Context, bot *tele.Bot) error {
	// Создаем HTTP клиента для загрузки изображения
	client := http.Client{
		Transport: cloneTransport(),
	}

	// Отправляем GET запрос для загрузки изображения
	res, err := client.Get(fileURL)
	if err != nil {
		log.Printf("Failed to fetch image: %s", err)
		return err
	}
	defer res.Body.Close()

	// Проверяем успешность запроса
	if res.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch image: %s", res.Status)
		return fmt.Errorf("failed to fetch image: %s", res.Status)
	}

	// Конвертируем изображение
	buffer := new(bytes.Buffer)
	err = convert(buffer, res.Body)
	if err != nil {
		log.Printf("Failed to convert image: %s", err)
		return err
	}

	// Отправляем сконвертированное изображение
	_, err = bot.Send(c.Chat(), &tele.Photo{
		File:    tele.FromReader(buffer),
		Caption: "Конвертированное изображение",
	})
	if err != nil {
		log.Printf("Failed to send image: %s", err)
		return err
	}

	return nil
}

func cloneTransport() *http.Transport {
	// Клонируем транспорт из DefaultTransport и настраиваем его
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true, // Этот параметр позволяет игнорировать проверку сертификата TLS (для примера)
		// Другие параметры безопасности, если необходимо
	}

	return transport
}

func convert(w io.Writer, r io.Reader) error {
	srcImage, err := imgconv.Decode(r)
	if err != nil {
		return fmt.Errorf("decode image: %w", err)
	}

	// Создаем экземпляр структуры imgconv.FormatOption
	formatOption := imgconv.FormatOption{Format: imgconv.PNG}

	err = imgconv.Write(w, srcImage, &formatOption)
	if err != nil {
		return fmt.Errorf("encode image: %w", err)
	}

	return nil
}
