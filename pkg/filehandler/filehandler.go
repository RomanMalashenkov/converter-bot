package filehandler

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RomanMalashenkov/tg_bot/pkg/converter"
	"github.com/RomanMalashenkov/tg_bot/pkg/httpclient"
	"github.com/sunshineplan/imgconv"
	tele "gopkg.in/telebot.v3"
)

func ConvertAndSendImage(fileURL string, c tele.Context, bot *tele.Bot, format imgconv.Format) error {
	// Создаем HTTP клиента для загрузки изображения
	client := http.Client{
		Transport: httpclient.CloneTransport(),
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
	err = converter.Convert(buffer, res.Body, format)
	if err != nil {
		log.Printf("Failed to convert image: %s", err)
		return err
	}

	// Создаем временный файл для сохранения сконвертированного изображения
	tempFileName := fmt.Sprintf("converted_image.%s", format)
	tempFile, err := os.Create(tempFileName)
	if err != nil {
		log.Printf("Failed to create temp file: %s", err)
		return err
	}
	defer tempFile.Close()

	// Записываем сконвертированное изображение во временный файл
	_, err = buffer.WriteTo(tempFile)
	if err != nil {
		log.Printf("Failed to write to temp file: %s", err)
		return err
	}

	// Отправляем сконвертированное изображение как файл
	fileToSend := &tele.Document{File: tele.FromDisk(tempFileName)}
	_, err = bot.Send(c.Chat(), fileToSend)
	if err != nil {
		log.Printf("Failed to send image file: %s", err)
		return err
	}

	// Удаляем временный файл после отправки
	tempFile.Close()
	err = os.Remove(tempFileName)
	if err != nil {
		log.Printf("Failed to delete temp file: %s", err)
	}

	return nil

}
