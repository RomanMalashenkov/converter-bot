package bot

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/RomanMalashenkov/tg_bot/pkg/api"
	"github.com/RomanMalashenkov/tg_bot/pkg/config"
	"github.com/RomanMalashenkov/tg_bot/pkg/converter"
	"github.com/RomanMalashenkov/tg_bot/pkg/filehandler"
	tele "gopkg.in/telebot.v3"
)

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
		fileURL, err := api.GetFileURL(b, doc.FileID)
		if err != nil {
			log.Printf("Failed to get file URL: %s", err)
			return err
		}
		// Отправляем запрос на конвертацию изображения
		err = filehandler.ConvertAndSendImage(fileURL, c, b)
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
