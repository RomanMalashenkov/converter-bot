package bot

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/RomanMalashenkov/config"
	"github.com/RomanMalashenkov/filehandler"

	tele "gopkg.in/telebot.v3"
)

func StartBot() {
	botConf, confErr := config.GetConfig()
	if confErr != nil {
		log.Fatal("No config")
	}

	filehandler.CheckFolder(botConf.config.Store)
	filehandler.CheckFolder(botConf.config.Store + "webm")
	filehandler.CheckFolder(botConf.config.Store + "mp4")

	b, err := tele.NewBot(tele.Settings{
		Token:  botConf.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Бот готов к работе")

	// ответ на команду /привет
	b.Handle("/start", func(c tele.Context) error {
		log.Printf("[LOG]: User: %s | Controller: /start ", c.Message().Sender.Username)
		_, err := b.Send(c.Chat(), "Привет, я бот-конвертер")
		return err
	})

	// ответ на команду /помощь
	b.Handle("/help", func(c tele.Context) error {
		log.Printf("[LOG]: User: %s | Controller: /help ", c.Message().Sender.Username)
		_, err := b.Send(c.Chat(), "Пришлите мне webm, а я вам - mp4") ////////////////////////////изменить ссобщ
		return err
	})

	//когда пользователь отправил док в чат
	b.Handle(tele.OnDocument, func(c tele.Context) error {
		log.Printf("[LOG]: User: %s | Controller: OnDocument ", c.Message().Sender.Username)
		doc := c.Message().Document

		//1111111проверка типа док
		if doc.MIME != "video/webm" {
			b.Send(c.Chat(), "Пожалуйста, пришлите мне файл формата webm для конвертации")
			return nil
		}

		//111111создаем пути для сохранения webm и mp4 файлов
		webmFilename := botConf.Store + "\\webm\\" + doc.FileID + ".webm"
		mp4Filename := botConf.Store + "\\mp4\\" + doc.FileID + ".mp4"
		messageFilename := strings.TrimSuffix(doc.FileName, ".webm") + ".mp4"

		b.Send(c.Chat(), "Загрузка...")

		if err := b.Download(&doc.File, webmFilename); err != nil {
			b.Send(c.Chat(), "Ошибка при загрузке файла")
			return nil
		}

		b.Send(c.Chat(), "Конвертация файла...")
		/////////////////////////Конвертация webm в mp4
		ffErr := filehandler.WebmToMp4(webmFilename, mp4Filename)
		if ffErr != nil {
			b.Send(c.Chat(), "Внутренняя ошибка сервера")
			log.Printf("webm: %s, mp4: %s", webmFilename, mp4Filename)
			log.Fatalf("Ошибка FFmpeg: %s", ffErr)
		}

		b.Send(c.Chat(), "Готово")
		mp4 := &tele.Video{File: tele.FromDisk(mp4Filename), FileName: messageFilename} //не Video!!!!

		b.Send(c.Chat(), mp4)
		//удалили файлы
		os.Remove(webmFilename)
		os.Remove(mp4Filename)

		return nil
	})

	b.Start()

}
