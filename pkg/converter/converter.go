package converter

import (
	"fmt"
	"io"
	"strings"

	"github.com/sunshineplan/imgconv"
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

func Convert(w io.Writer, r io.Reader) error {
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
