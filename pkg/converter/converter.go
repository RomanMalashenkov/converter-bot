package converter

import (
	"fmt"
	"io"
	"log"
	"strings"

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
			log.Printf("Уникальный ID кнопки: %s, Текст кнопки: %s\n", btn.Data, btn.Text)
		}
		inlineBtns = append(inlineBtns, inlineRow)
	}
	return inlineBtns
}

// функция конвертации
func Convert(w io.Writer, r io.Reader, format imgconv.Format) error {
	srcImage, err := imgconv.Decode(r)
	if err != nil {
		return fmt.Errorf("decode image: %w", err)
	}

	// Создаем экземпляр структуры imgconv.FormatOption
	formatOption := imgconv.FormatOption{Format: format}

	err = imgconv.Write(w, srcImage, &formatOption)
	if err != nil {
		return fmt.Errorf("encode image: %w", err)
	}

	return nil
}
