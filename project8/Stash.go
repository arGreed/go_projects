package stash

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// Информация о продукте.
type Product struct {
	Id     int    `json:"Id"`
	Name   string `json:"Name"`
	Amount int    `json:"Amount"`
	Sold   int    `json:"Sold"`
	Price  int    `json:"Price"`
}

type usrStash struct {
	id     int
	amount int
	price  int
}

var (
	products    = make(map[int]Product)
	stash       []usrStash
	storagePath string = "project8/storage.json"
	curId       int    = 1
)

var (
	errInvalidInput = "Некорректный ввод, повторите попытку!"
)

// Команды администратора.
const (
	admCommandAdd           = "add"
	admCommandDelete        = "delete"
	admCommandRename        = "rename"
	admCommandCalculate     = "acalculate"
	admCommandCalculateSold = "acalcullatesold"
	admCommandShow1         = "ashow1"
	admCommandShowAll       = "ashow"
	admCommandChangeAmount  = "amountChange"
)

// Общие команды.
const (
	comCommandExit = "exit"
)

// Команды пользователя
const (
	usrCommandShow      = "show"
	usrCommandAddToCart = "add"
	usrCommandBuy       = "buy"
)

func userAdd() error {
	var s string
	var idea usrStash
	var flag bool = true
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Введите идентификаторы заинтересовавшего вас товара:")
	for {
		if scanner.Scan() {
			s = scanner.Text()
		}
		if s == comCommandExit {
			return nil
		}
		id, err := strconv.Atoi(s)

		if err == nil && products[id].Id != 0 {
			for {
				fmt.Print("Введите количество покупаемого товара:")
				if scanner.Scan() {
					s = scanner.Text()
				}
				if s == comCommandExit {
					return nil
				}
				cnt, err := strconv.Atoi(s)
				if err == nil && products[id].Amount-cnt >= 0 {
					idea.id = id
					idea.amount = cnt
					idea.price = products[id].Price * cnt
					stash = append(stash, idea)
					fmt.Println("Хотите ли добавить ещё что-то в корзину ещё что-то?")
					if scanner.Scan() {
						s = scanner.Text()
						if s != "да" {
							return nil
						}
					}
				}
				if flag {
					fmt.Println("Введено некорректное значение, повторите ввод!")
				}
			}
		}
		if flag {
			fmt.Println("Введено некорректное значение, повторите ввод!")
		}
		flag = true
	}
}

func userShow() error {
	if len(products) > 0 {
		for _, i := range products {
			fmt.Println("name:", i.Name, " price:", i.Price, " id:", i.Id)
		}
	} else {
		fmt.Printf("Нет информации о товарах.")
	}
	return nil
}

func userBuy() {
	var money int = 0
	var product Product
	if len(stash) == 0 {
		fmt.Println("Сперва добавьте товар в корзину!")
		return
	}
	for _, i := range stash {
		money += i.price
		product = products[i.id]
		product.Sold += i.amount
		product.Amount -= i.amount
		products[i.id] = product
	}
	save()
	clear(stash)
	log.Println(time.Now(), " Совершена покупка на ", money, " Валюты")
}

// Суперпользователь
func adminMode() {
	var s string
	fmt.Println("Список доступных команд:")
	fmt.Println("1)", admCommandAdd)
	fmt.Println("2)", admCommandDelete)
	fmt.Println("3)", admCommandRename)
	fmt.Println("4)", admCommandCalculate)
	fmt.Println("5)", admCommandCalculateSold)
	fmt.Println("6)", admCommandShow1)
	fmt.Println("7)", admCommandShowAll)
	fmt.Println("8)", admCommandChangeAmount)
	fmt.Println("9)", comCommandExit)
	for {
		fmt.Print("Команда:")
		fmt.Scanln(&s)
		switch s {
		case admCommandAdd:
			err := admAdd()
			if err != nil {
				log.Printf(err.Error())
				return
			}
			s = ""
		case admCommandDelete:
			err := admDel()
			if err != nil {
				log.Printf(err.Error())
				return
			}
			s = ""
		case admCommandRename:
			err := admRename()
			if err != nil {
				log.Printf(err.Error())
				return
			}
			s = ""
		case admCommandCalculate:
			err := admCalculate()
			if err != nil {
				log.Printf(err.Error())
				return
			}
			s = ""
		case admCommandCalculateSold:
			err := admCalculatesold()
			if err != nil {
				log.Printf(err.Error())
				return
			}
			s = ""
		case admCommandShow1:
			err := admShowProduct()
			if err != nil {
				log.Printf(err.Error())
				return
			}
			s = ""
		case admCommandShowAll:
			err := admShowProducts()
			if err != nil {
				log.Printf(err.Error())
				return
			}
			s = ""
		case admCommandChangeAmount:
			err := admChangeAmount()
			if err != nil {
				log.Printf(err.Error())
				return
			}
			s = ""
		case comCommandExit:
			return
		default:
			fmt.Println("Введена некорректная команда, повторите ввод!")
			s = ""
		}
	}
}

// Покупатель
func userMode() {
	var s string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Список доступных пользователю команд:")
	fmt.Println("1)", usrCommandShow)
	fmt.Println("2)", usrCommandAddToCart)
	fmt.Println("3)", usrCommandBuy)
	fmt.Println("4)", comCommandExit)
	fmt.Print("Введите команду:")
	for {
		if scanner.Scan() {
			s = scanner.Text()
			switch s {
			case usrCommandShow:
				err := userShow()
				if err != nil {
					log.Println(time.Now(), err.Error())
				}
				s = ""
			case usrCommandAddToCart:
				err := userAdd()
				if err != nil {
					log.Println(time.Now(), err.Error())
				}
				s = ""
			case usrCommandBuy:
				userBuy()
				s = ""
				return
			case comCommandExit:
				return
			}
		}
	}
}

func storageInit() error {
	if fileinfo, _ := os.Stat(storagePath); fileinfo.Size() == 0 {
		return nil
	}

	file, err := os.Open(storagePath)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&products)

	if err != nil {
		return err
	}

	for _, i := range products {
		curId = i.Id + 1
	}

	return nil
}

func save() error {
	file, err := os.OpenFile(storagePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(&products)
	// Если сверху выпала ошибка, вернёт ошибку, если нет - nil
	return err
}

func Shop() {
	var mode string
	if err := storageInit(); err != nil {
		log.Println(err)
		return
	}
	fmt.Scanln(&mode)
	if mode == "admin" {
		adminMode()
	} else {
		userMode()
	}
}

func admAdd() error {
	var s string
	var buf int = 0
	var product Product
	var err error
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Введите наименование товара:")
		if scanner.Scan() {
			s = scanner.Text()
		}
		if s == comCommandExit {
			log.Println(time.Now(), " Отмена добавления задачи.")
			return nil
		}
		if s != "" {
			break
		} else {
			fmt.Println(errInvalidInput)
			log.Println(time.Now(), " При добавлении товара введено некорректное наименование!")
		}
	}
	product.Name = s
	for {
		fmt.Println("Введите количество товара на складе:")
		if scanner.Scan() {
			s = scanner.Text()
		}
		if s == comCommandExit {
			log.Println(time.Now(), " Отмена добавления задачи.")
			return nil
		}
		buf, err = strconv.Atoi(s)
		if err == nil && buf > 0 {
			break
		}
		fmt.Println(errInvalidInput)
		log.Println(time.Now(), " При добавлении товара введено некорректное количество!")
	}
	product.Amount = buf
	for {
		fmt.Println("Введите цену за единицу товара:")
		if scanner.Scan() {
			s = scanner.Text()
		}
		if s == comCommandExit {
			log.Println(time.Now(), " Отмена добавления задачи.")
			return nil
		}
		buf, err = strconv.Atoi(s)
		if err == nil && buf > 0 {
			break
		}
		fmt.Println(errInvalidInput)
		log.Println(time.Now(), " При добавлении товара введена некорректная цена!")
	}
	product.Price = buf

	product.Id = curId
	products[curId] = product
	err = save()
	if err != nil {
		return err
	}
	curId++
	return nil
}

func admDel() error {
	scanner := bufio.NewScanner(os.Stdin)
	var buf string
	fmt.Print("Введите идентификатор удаляемого товара:")
	for {
		if scanner.Scan() {
			buf = scanner.Text()
		}
		if buf == comCommandExit {
			log.Println(time.Now(), " Отмена добавления задачи.")
			return nil
		}
		id, err := strconv.Atoi(buf)
		if err != nil {
			log.Println(time.Now(), " ", errInvalidInput)
			return err
		} else {
			if products[id].Id == 0 {
				log.Println(time.Now(), " Попытка удалить несуществующий элемент")
			} else {
				delete(products, id)
				log.Println(time.Now(), " Элемент успешно удалён.")
				save()
				return nil
			}
		}
	}
}

func admRename() error {
	var s string
	var id int = 0
	var err error
	var product Product
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Введите идентификатор переименовываемого товара:")
	for {
		if scanner.Scan() {
			s = scanner.Text()
		}
		if s == comCommandExit {
			log.Println(time.Now(), " Отмена добавления задачи.")
			return nil
		}

		id, err = strconv.Atoi(s)

		if products[id].Id != 0 && err == nil && s != "" {
			break
		}

		fmt.Println("Элемент не найден или введено некорректное значение, повторите ввод!")
	}
	product = products[id]
	fmt.Print("Введите новое наименование:")
	for {
		if scanner.Scan() {
			s = scanner.Text()
		}
		if s == comCommandExit {
			log.Println(time.Now(), " Отмена добавления задачи.")
			return nil
		}
		if s != "" {
			product.Name = s
			products[id] = product
			save()
			return nil
		}
		fmt.Println("Введено некорректное значение, повторите ввод!")
	}
}

func admCalculate() error {
	var s string
	var id int
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Ведите идентификатор товара:")
	for {
		if scanner.Scan() {
			s = scanner.Text()
		}
		if s == comCommandExit {
			log.Println(time.Now(), " Отмена добавления задачи.")
			return nil
		}

		id, err = strconv.Atoi(s)

		if products[id].Id != 0 && err == nil && s != "" {
			break
		}

		fmt.Println("Элемент не найден или введено некорректное значение, повторите ввод!")
	}

	fmt.Println("Суммарная стоимость товара ", products[id].Name, " на складе составляет ", products[id].Price*products[id].Amount)
	return nil
}

func admCalculatesold() error {
	var s string
	var id int
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Ведите идентификатор товара:")
	for {
		if scanner.Scan() {
			s = scanner.Text()
		}
		if s == comCommandExit {
			log.Println(time.Now(), " Отмена добавления задачи.")
			return nil
		}

		id, err = strconv.Atoi(s)

		if products[id].Id != 0 && err == nil && s != "" {
			break
		}

		fmt.Println("Элемент не найден или введено некорректное значение, повторите ввод!")
	}

	fmt.Println("Суммарная стоимость проданного товара ", products[id].Name, " составляет ", products[id].Sold*products[id].Price)
	return nil
}

func admChangeAmount() error {
	var s string
	var id, cnt int
	var err error
	var product Product
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Ведите идентификатор товара:")
	for {
		if scanner.Scan() {
			s = scanner.Text()
		}
		if s == comCommandExit {
			log.Println(time.Now(), " Отмена добавления задачи.")
			return nil
		}

		id, err = strconv.Atoi(s)

		if products[id].Id != 0 && err == nil && s != "" {
			break
		}

		fmt.Println("Элемент не найден или введено некорректное значение, повторите ввод!")
	}
	product = products[id]
	fmt.Println("Введите количество полученного или списанного товара:")
	for {
		if scanner.Scan() {
			s = scanner.Text()
		}
		if s == comCommandExit {
			log.Println(time.Now(), " Отмена добавления задачи.")
			return nil
		}

		cnt, err = strconv.Atoi(s)

		if err == nil && (cnt > 0 || products[id].Amount-cnt > 0) {
			product.Amount += cnt
			products[id] = product
			save()
			return nil
		}

		fmt.Println("Введено некорректное значение, повторите ввод!")
	}
}

func admShowProduct() error {
	var s string
	var id int
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Ведите идентификатор товара:")
	for {
		if scanner.Scan() {
			s = scanner.Text()
		}
		if s == comCommandExit {
			log.Println(time.Now(), " Отмена добавления задачи.")
			return nil
		}

		id, err = strconv.Atoi(s)

		if products[id].Id != 0 && err == nil && s != "" {
			break
		}

		fmt.Println("Элемент не найден или введено некорректное значение, повторите ввод!")
	}

	fmt.Println("id:", products[id].Id, " name:", products[id].Name, " amount:", products[id].Amount, " price:", products[id].Price, "sold:", products[id].Sold)
	return nil
}

func admShowProducts() error {
	if len(products) > 0 {
		for _, i := range products {
			fmt.Println("id:", i.Id, " name:", i.Name, " amount:", i.Amount, " price:", i.Price, "sold:", i.Sold)
		}
	} else {
		fmt.Printf("Нет информации о отоварах.")
	}
	return nil
}
