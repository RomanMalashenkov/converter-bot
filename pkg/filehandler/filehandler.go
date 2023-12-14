package filehandler

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RomanMalashenkov/tg_bot/pkg/converter"
	"github.com/RomanMalashenkov/tg_bot/pkg/httpclient"
	tele "gopkg.in/telebot.v3"
)

// проверка существования папок и их создания (если отстутствуют)
func CheckFolder(store string) {
	if _, err := os.Stat(store); os.IsNotExist(err) {
		log.Printf("Добавлена папка %s", store)
		err = os.Mkdir(store, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ConvertAndSendImage(fileURL string, c tele.Context, bot *tele.Bot) error {
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
	err = converter.Convert(buffer, res.Body)
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
