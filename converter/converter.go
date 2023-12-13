package converter

import (
	"log"

	"github.com/sunshineplan/imgconv"
)

type Ras struct {
	FormatFile imgconv.Format
}

// узнать формат файла
func FormatFile() error {
	ras, err := imgconv.FormatFromExtension(FormatFile)
	if ras !=  {
		log.Printf("Пожалуйста, пришлите мне файл формата webm для конвертации")

	}
}

type FormatOption struct {
	Format       Format
	EncodeOption []EncodeOption
}
