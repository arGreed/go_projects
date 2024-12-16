package notes

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
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

	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&noteLst)

	if err != nil {
		log.Println(err)
		return err
	}
	for _, i := range *noteLst {
		*id = i.Id + 1
	}

	return nil
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

func getInput(arg string) string {
	scanner := bufio.NewScanner(os.Stdin)
	var s string
	fmt.Print(arg)
	if scanner.Scan() {
		s = scanner.Text()
	}
	return s
}

func addNote(allNotes *map[int]Note, lastId *int) error {
	var s string
	var note Note
	for {
		s = getInput("Введите заголовок новой заметки:")

		if s == commandExit {
			return nil
		}
		if s != "" {
			note.Name = s
			break
		}
		fmt.Println("Наименование не должно быть пустым, повторите ввод!")
	}
	for {
		s = getInput("Введите описание заметки:")
		if s == commandExit {
			return nil
		}
		if s != "" {
			note.Description = s
			break
		}
		fmt.Println("Описание не должно быть пустым, повторите ввод!")
	}
	for {
		s = getInput("Оцените важность этой заметки (чем больше число, тем менее важная заметка)")
		if s == commandExit {
			return nil
		}
		buf, err := strconv.Atoi(s)
		if err == nil && buf > 0 {
			note.StatusRate = buf
			break
		}
		fmt.Println("Введено некорректное значение, повторите ввод")
	}
	note.Id = *lastId + 1
	*lastId++
	(*allNotes)[note.Id] = note
	err := storageSave(allNotes)
	return err
}

func deleteNote(allNotes *map[int]Note, lastId *int) error {
	var s string
	var buf int
	var err error
	for {
		s = getInput("Введите идентификатор удаляемой заметки")
		if s == commandExit {
			return nil
		}
		buf, err = strconv.Atoi(s)
		if err == nil && buf > 0 && buf <= *lastId {
			break
		}
	}
	if (*allNotes)[buf].Name != "" {
		delete(*allNotes, buf)
	}
	storageSave(allNotes)
	return nil
}

func updateNote(allNotes *map[int]Note, lastId *int) error {
	var s string
	var buf int
	var err error
	var note Note
	if len(*allNotes) == 0 {
		fmt.Println("Хранилище пусто!")
		return nil
	}
	for {
		s = getInput("Введите идентификатор изменяемой заметки")
		if s == commandExit {
			return nil
		}
		buf, err = strconv.Atoi(s)
		if err == nil && buf > 0 && buf <= *lastId && (*allNotes)[buf].Name != "" {
			break
		}
		fmt.Println("Некорректный ввод или введён несуществующий идентификатор, повторите попытку!")
	}
	note.Id = buf
	for {
		s = getInput("Введите заголовок новой заметки:")

		if s == commandExit {
			return nil
		}
		if s != "" {
			note.Name = s
			break
		}
		fmt.Println("Наименование не должно быть пустым, повторите ввод!")
	}
	for {
		s = getInput("Введите описание заметки:")
		if s == commandExit {
			return nil
		}
		if s != "" {
			note.Description = s
			break
		}
		fmt.Println("Описание не должно быть пустым, повторите ввод!")
	}
	for {
		s = getInput("Оцените важность этой заметки (чем больше число, тем менее важная заметка):")
		if s == commandExit {
			return nil
		}
		buf, err := strconv.Atoi(s)
		if err == nil && buf > 0 {
			note.StatusRate = buf
			break
		}

	}
	(*allNotes)[note.Id] = note
	storageSave(allNotes)
	return nil
}

func reRateNote(allNotes *map[int]Note, lastId *int) error {
	var s string
	var buf, buf1 int
	var err error
	var note Note
	if len(*allNotes) == 0 {
		fmt.Println("Хранилище пусто!")
		return nil
	}
	for {
		s = getInput("Введите идентификатор изменяемой заметки:")
		if s == commandExit {
			return nil
		}
		buf, err = strconv.Atoi(s)
		if err == nil && buf > 0 && buf <= *lastId && (*allNotes)[buf].Name != "" {
			break
		}
		fmt.Println("Некорректный ввод или введён несуществующий идентификатор, повторите попытку!")
	}
	note = (*allNotes)[buf]
	for {
		s = getInput("Введите новый рейтинг заметки:")
		if s == commandExit {
			return nil
		}
		buf1, err = strconv.Atoi(s)
		if err == nil && buf1 > 0 {
			break
		}
		fmt.Println("Некорректный ввод или введён несуществующий идентификатор, повторите попытку!")
	}
	note.StatusRate = buf1
	(*allNotes)[buf] = note
	storageSave(allNotes)
	return nil
}

func noteShow(allNotes *map[int]Note) {
	var s string
	var buf int
	var err error
	for {
		s = getInput("Введите идентификатор интересующей заметки:")
		if s == commandExit {
			return
		}
		buf, err = strconv.Atoi(s)
		if err == nil && buf > 0 && (*allNotes)[buf].Name != "" {
			break
		}
		fmt.Println("Некорректный ввод или введён несуществующий идентификатор, повторите попытку!")
	}
	fmt.Println("id:", (*allNotes)[buf].Id, " description:", (*allNotes)[buf].Description, " name:", (*allNotes)[buf].Name, " rating:", (*allNotes)[buf].StatusRate)
}

func noteAll(allNotes *map[int]Note) {
	var list []Note = make([]Note, 0, len(*allNotes))
	var buf Note
	for _, i := range *allNotes {
		list = append(list, i)
	}
	for i := 0; i < len(list); i++ {
		for j := i + 1; j < len(list); j++ {
			if list[i].StatusRate > list[j].StatusRate {
				buf = list[i]
				list[i] = list[j]
				list[j] = buf
			}
		}
	}
	for _, i := range list {
		fmt.Println("id:", i.Id, " description:", i.Description, " name:", i.Name, " rating:", i.StatusRate)
	}
}

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
			err := addNote(allNotes, lastId)
			if err != nil {
				return err
			}
			s = ""
		case commandDelete:
			err := deleteNote(allNotes, lastId)
			if err != nil {
				return err
			}
			s = ""
		case commandUpdate:
			err := updateNote(allNotes, lastId)
			if err != nil {
				return err
			}
			s = ""
		case commandChangeImp:
			err := reRateNote(allNotes, lastId)
			if err != nil {
				return err
			}
			s = ""
		case commandShow:
			noteShow(allNotes)
		case commandShowAll:
			noteAll(allNotes)
		case commandExit:
			return nil
		default:
			fmt.Println("Введена недопустимая задача, повторите попытку снова!")
		}
	}
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
