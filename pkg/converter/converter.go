package converter

import (
	"strings"
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
