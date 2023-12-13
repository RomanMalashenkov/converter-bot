package filehandler

import (
	"log"
	"os"
	"os/exec"
)

type File struct {
	store string
}

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

// ф-я для конвертации webm в mp4
func WebmToMp4(in string, out string) error {
	cmd := exec.Command("ffmpeg", "-i", in, out)

	log.Print(cmd.Args) //лог команду, которая будет выполнена

	err := cmd.Run()

	if err != nil {
		log.Print(cmd.Stderr) //Логаем вывод ошибок
		log.Print(cmd.Stdout) //Логает стандартный вывод

		return err
	}
	return nil
}
