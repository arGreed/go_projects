package notes

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const logFile = "project10/test.log"
const storageFile = "project10/storage.json"

type Note struct {
	Id          int    `json:"Id"`
	StatusRate  int    `json:"StatusRate"`
	Name        string `json:"Name"`
	Description string `json:""`
}

func logPrepare() (*os.File, error) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.SetOutput(file)
	return file, err
}

func storageInit(noteLst *map[int]Note, id *int) error {
	if fileInfo, _ := os.Stat(storageFile); fileInfo.Size() == 0 {
		log.Println("Хранилище пусто.")
		return nil
	}

	file, err := os.Open(storageFile)

	if err != nil {
		log.Println(err)
		return err
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&noteLst)

	defer file.Close()

	for _, i := range *noteLst {
		*id = i.Id + 1
	}

	return err
}

// Сохранение изменений списка в хранилище.
func storageSave(noteLst *map[int]Note) error {
	file, err := os.OpenFile(storageFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	encoder := json.NewEncoder(file)

	encoder.Encode(noteLst)

	return nil
}

const (
	commandAdd       = "add"
	commandDelete    = "delete"
	commandUpdate    = "update"
	commandChangeImp = "imp"
	commandShow      = "see"
	commandShowAll   = "show"
	commandExit      = "exit"
)

func operateNotes(allNotes *map[int]Note, lastId *int) error {
	var s string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Введите одну из доступных команд:")
	fmt.Println(commandAdd, "- добавление заметки.")
	fmt.Println(commandDelete, "- удаление заметки")
	fmt.Println(commandUpdate, "- изменение имени и описания заметки.")
	fmt.Println(commandChangeImp, "- изменение приоритета заметки.")
	fmt.Println(commandShow, "- увидеть заметку по заданному идентификатору.")
	fmt.Println(commandShowAll, "- увидеть список всех заметок.")
	fmt.Println(commandExit, "- выход из приложения.")
	for {
		if scanner.Scan() {
			s = scanner.Text()
		}
		switch s {
		case commandAdd:
		case commandDelete:
		case commandUpdate:
		case commandChangeImp:
		case commandShow:
		case commandShowAll:
		case commandExit:
		default:
			fmt.Println("Введена недопустимая задача, повторите попытку снова!")
		}
	}
	return nil
}

func Notes() {
	// Подготовка лог файла.
	file, err := logPrepare()
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	// В других проектах сделал через глобальные переменные, сейчас переделаю без них.
	var noteLst = make(map[int]Note)
	var noteCurId int = 1

	// Инициализация хранилища.
	err = storageInit(&noteLst, &noteCurId)

	if err != nil {
		log.Println(err)
		return
	}

	err = operateNotes(&noteLst, &noteCurId)

	if err != nil {
		log.Println(err)
		return
	}
}
