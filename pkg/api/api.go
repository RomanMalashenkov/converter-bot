package api

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func GetFileURL(bot *tele.Bot, fileID string) (string, error) {
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

/*
// Отправка сообщения с изображением через API Telegram
func SendImage(bot *tele.Bot, chat tele.Recipient, image tele.Photo) error {
	_, err := bot.Send(chat, &image)
	if err != nil {
		log.Printf("Failed to send image: %s", err)
		return err
	}

	return nil
}
*/
