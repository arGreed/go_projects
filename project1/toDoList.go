package toDoList

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Структура задачи.
type Task struct {
	// Техническое поле, идентификатор задачи.
	Id int `json:"Id"`
	// Наименование задачи.
	Name string `json:"Name"`
	// Описание задачи.
	Description string `json:"Description"`
	// Статус задачи.
	Status bool `json:"Status"`
}

// Список команд.
const (
	CommandAdd    = "add"
	CommandDelete = "delete"
	CommandUpdate = "update"
	CommandShow   = "show"
	CommandAll    = "all"
	CommandExit   = "exit"
)

// Список пользовательских ошибок.
var (
	errWrongInput = errors.New("введено некорректное значение")
)

// Глобальные переменные.
var (
	// Переменная, используемая для проставления идентификатора задач.
	curId int = 1
	// Список задач, выгруженных в память.
	TaskList = make(map[int]Task)
	// Путь к хранилищу задач.
	storage = "project1/storage.json"
)

// Функция, используемая для выгрузки списка задач в оперативную память.
func taskInit() error {
	// Проверка на наличие записей в хранилище.
	fileInfo, _ := os.Stat(storage)
	if fileInfo.Size() == 0 {
		return nil
	}

	// Открытый файл для считывания списка задач.
	file, err := os.Open(storage)

	// Ошибка при открытии файла.
	if err != nil {
		return err
	}

	// Декодировщик json-файла.
	decoder := json.NewDecoder(file)

	// Инициализация списка задач данными из хранилища.
	err = decoder.Decode(&TaskList)

	// Проверка на успешность инициализации списка задач.
	if err != nil {
		return err
	}

	// Выборка максимального идентификатора + 1.
	for _, i := range TaskList {
		curId = i.Id + 1
	}

	// Если успешно проинициализировали список задач.
	return nil
}

// Функция сохранения данных в хранилище.
func save() error {
	// Открытие файла с определёнными флагами.
	file, err := os.OpenFile(storage, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	// Если произошла ошибка при открытии файла, функция закончит выполнение, иначе будет вызов отложенного закрытия файла.
	if err != nil {
		return err
	} else {
		// Отложенное закрытие файла.
		defer file.Close()
	}

	// Кодировщик данных в json-файл.
	encoder := json.NewEncoder(file)

	// Занесение списка задач в хранилище.
	err = encoder.Encode(TaskList)

	// Если столкнулись с ошибкой записи данных в файл.
	if err != nil {
		return err
	}

	// Если успешно записали список задач в хранилище.
	return nil
}

// Функция добавления задачи в список задач.
func addTask() error {
	// Считываемая задача.
	var task Task

	// Консольный ввод наименования и описания задачи.
	fmt.Print("Введите наименование задачи:")

	// Считывание.
	fmt.Scanln(&task.Name)
	fmt.Print("Введите Описание задачи:")

	// Считывание.
	fmt.Scanln(&task.Description)

	// Проверка на корректность ввода.
	if task.Name == "" || task.Description == "" {
		return errWrongInput
	}

	// Присваивание задаче нового идентификатора.
	task.Id = curId

	// Подготовка идентификатора для следующей задачи.
	curId++

	// Запись задачи в список задач.
	// Статус задачи по дефолту имеет значение false
	TaskList[task.Id] = task

	// Фиксация изменений.
	err := save()

	// Ошибка при фиксации изменений.
	if err != nil {
		return err
	}
	// Если успешно добавили задачу.
	return nil
}

// Функция удаления задачи.
func delTask() error {
	// Считываемый идентификатор задачи.
	var id int

	// Консольный ввод идентификатора задачи.
	fmt.Print("Введите идентификатор задачи: ")
	_, err := fmt.Scanln(&id)

	if err != nil {
		return err
	}

	// Проверка существования удаляемого элемента.
	if TaskList[id].Id == 0 {
		return errors.New("запрошен несуществующий элемент")
	}

	// Удаление элемента из списка задач.
	delete(TaskList, id)

	// Фиксация изменений в хранилище.
	err = save()

	// Оповещение о успешном удалении задачи.
	if err == nil {
		fmt.Println("Задача успешно удалена!")
	} else {
		return err
	}

	// Если задача успешно удалена.
	return nil
}

// Функция изменения задачи.
func updTask() error {
	// Считываемый идентификатор задачи.
	var id int

	// Обновлённые параметры задачи.
	var task Task

	// Блок считывания идентификатора задачи.
	fmt.Print("Введите идентификатор задачи: ")
	_, err := fmt.Scanln(&id)

	// Проверка корректности ввода.
	if err != nil {
		return err
	}

	// Проверка существования элемента в хранилище.
	if TaskList[id].Id == 0 {
		return errors.New("запрошен несуществующий элемент")
	}

	// Блок считывания обновлённых параметров задачи.
	fmt.Println("Введите обновлённое наименование задачи: ")
	fmt.Scanln(&task.Name)
	fmt.Println("Введите обновлённое Описание задачи: ")
	fmt.Scanln(&task.Description)

	// Подготовка к записи в хранилище.
	task.Id = id

	// Запись в хранилище.
	TaskList[id] = task

	// Фиксация изменений в хранилище.
	err = save()

	// Проверка успешности фиксаций изменений.
	if err == nil {
		fmt.Printf("Структура задачи была успешно обновлена!")
		return nil
	} else {
		return err
	}
}

// Функция, отображающая задачу по конкретному идентификатору.
func showTask() error {
	// Идентификатор отображаемой задачи.
	var id int

	// Блок считывания идентификатора задачи.
	fmt.Print("Введите идентификатор задачи: ")
	_, err := fmt.Scanln(&id)

	// Проверка корректности ввода.
	if err != nil {
		return err
	}

	// Проверка существования элемента в хранилище.
	if TaskList[id].Id == 0 {
		return errors.New("запрошен несуществующий элемент")
	}

	// Вывод запрошенного элемента.
	fmt.Println("Id: ", TaskList[id].Id, " Name: ", TaskList[id].Name, " Description: ", TaskList[id].Description, " Status: ", TaskList[id].Status)

	// Успешно вывели элемент.
	return nil
}

// Отображает весь список задач.
func showAllTask() {
	for _, i := range TaskList {
		fmt.Println("Id: ", i.Id, " Name: ", i.Name, " Description: ", i.Description, " Status: ", i.Status)
	}
	if len(TaskList) == 0 {
		fmt.Println("Список пуст, скорее добавьте свою первую задачу!")
	}
}

// Маршрутизатор приложения.
func operateTasks() error {
	// Вводимая команда.
	var s string = ""

	// Ошибка исполнения программы.
	var err error

	// Проверка успешности загрузки данных из хранилища.
	if err = taskInit(); err != nil {
		return err
	}

	for {
		// Считывание команд.
		fmt.Print("Введите команду, которую вы собираетесь выполнить (add, delete, update, show, all, exit):")
		fmt.Scanln(&s)
		switch s {
		// Добавление задачи.
		case CommandAdd:
			if err = addTask(); err != nil {
				fmt.Println(err)
			}
			s = ""
		// Удаление задачи.
		case CommandDelete:
			if err = delTask(); err != nil {
				fmt.Println(err)
			}
			s = ""
		// Изменение задачи.
		case CommandUpdate:
			if err = updTask(); err != nil {
				fmt.Println(err)
			}
			s = ""
		// Отображение конкретной задачи.
		case CommandShow:
			if err = showTask(); err != nil {
				fmt.Println(err)
			}
			s = ""
		// Отображение списка задач.
		case CommandAll:
			showAllTask()
			s = ""
		// Завершение работы приложения.
		case CommandExit:
			return nil
		// При вводе некорректной задачи.
		default:
			fmt.Println("Введена некорректная задача, повторите ввод!")
			continue
		}
	}
}

// Запуск модуля.
func ToDOList() {
	if err := operateTasks(); err != nil {
		fmt.Println("При выполнении программы возникла ошибка: ", err)
	}
	defer func() {
		err := save()
		if err == nil {
			fmt.Println("Ваши задачи были успешно сохранены, хорошего дня!")
		} else {
			fmt.Println("Произошла ошибка при сохранении задач, изменения не были зафиксированы, хорошего дня!")
		}
	}()
}
