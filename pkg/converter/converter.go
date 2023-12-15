package converter

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/RomanMalashenkov/tg_bot/pkg/httpclient"
	"github.com/sunshineplan/imgconv"
	tele "gopkg.in/telebot.v3"
)

// проверяет, поддерживается ли формат для конвертации
func IsSupported(ext string) bool {
	supportedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".tif", ".tiff", ".bmp", ".pdf"}
	lowercaseExt := strings.ToLower(ext)

	for _, supportedExt := range supportedExts {
		if supportedExt == lowercaseExt {
			return true
		}
	}
	return false
}

func ConvertToInlineButtons(btns [][]tele.Btn) [][]tele.InlineButton {
	var inlineBtns [][]tele.InlineButton
	for _, btnRow := range btns {
		var inlineRow []tele.InlineButton
		for _, btn := range btnRow {
			inlineRow = append(inlineRow, tele.InlineButton{
				Unique: btn.Data,
				Text:   btn.Text,
			})
		}
		inlineBtns = append(inlineBtns, inlineRow)
	}
	return inlineBtns
}

func ConvertAndSendImage(fileURL string, c tele.Context, bot *tele.Bot, format imgconv.Format) error {
	// Создаем HTTP клиента для загрузки изображения
	client := http.Client{
		Transport: httpclient.CloneTransport(),
	}

	// Отправляем GET запрос для загрузки изображения
	res, err := client.Get(fileURL)
	if err != nil {
		log.Printf("Не удалось получить изображение: %s", err)
		return err
	}
	defer res.Body.Close()

	// Проверяем успешность запроса
	if res.StatusCode != http.StatusOK {
		log.Printf("Не удалось получить изображение...: %s", res.Status)
		return fmt.Errorf("не удалось получить изображение...: %s", res.Status)
	}

	// Конвертируем изображение
	buffer := new(bytes.Buffer)
	err = Convert(buffer, res.Body, format)
	if err != nil {
		log.Printf("Не удалось преобразовать изображение: %s", err)
		return err
	}

	// Создаем временный файл для сохранения сконвертированного изображения
	tempFileName := fmt.Sprintf("converted_image.%s", format)
	tempFile, err := os.Create(tempFileName)
	if err != nil {
		log.Printf("Не удалось создать временный файл: %s", err)
		return err
	}

	// Записываем сконвертированное изображение во временный файл
	_, err = buffer.WriteTo(tempFile)
	if err != nil {
		log.Printf("Не удалось записать во временный файл: %s", err)
		return err
	}
	////////////////////////////тут назва ф над добавить???
	// Отправляем сконвертированное изображение как файл
	fileToSend := &tele.Document{File: tele.FromDisk(tempFileName), FileName: tempFileName}
	_, err = bot.Send(c.Chat(), fileToSend)
	if err != nil {
		log.Printf("Не удалось отправить файл изображения: %s", err)
		return err
	}

	// Удаляем временный файл после отправки
	tempFile.Close()
	err = os.Remove(tempFileName)
	if err != nil {
		log.Printf("Не удалось удалить временный файл: %s", err)
	}

	return nil

}

// функция конвертации
func Convert(w io.Writer, r io.Reader, format imgconv.Format) error {
	srcImage, err := imgconv.Decode(r)
	if err != nil {
		return fmt.Errorf("не удалось декодировать изображение: %w", err)
	}

	// Создаем экземпляр структуры imgconv.FormatOption
	formatOption := imgconv.FormatOption{Format: format}

	err = imgconv.Write(w, srcImage, &formatOption)
	if err != nil {
		return fmt.Errorf("не удалось кодировать изображение: %w", err)
	}

	return nil
}
