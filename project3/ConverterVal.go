package project3

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/parnurzeal/gorequest"
)

type Coin struct {
	Id       string  `json:"ID"`
	NumCode  string  `json:"NumCode"`
	CharCode string  `json:"CharCode"`
	Name     string  `json:"Name"`
	Nominal  int     `json:"Nominal"`
	Value    float64 `json:"Value"`
	Previous float64 `json:"Previous"`
}

type ValueResponse struct {
	Value map[string]Coin `json:"Valute"`
}

var (
	coinList = make(map[string]Coin)
)

func getMoney() error {
	var url = "https://www.cbr-xml-daily.ru/daily_json.js"

	resp, _, err := gorequest.New().Get(url).End()

	if err != nil {
		resp.Body.Close()
		return errors.New("error getting money")
	}
	defer resp.Body.Close()

	var valueResponse ValueResponse

	err1 := json.NewDecoder(resp.Body).Decode(&valueResponse)

	if err1 != nil {
		fmt.Println(err)
		return errors.New("error decoding body")
	}

	for name, i := range valueResponse.Value {
		coinList[name] = i
	}

	return nil
}

func readConsole(first, second *string, amount *float64) {
	fmt.Print("Введите наименование конвертируемой валюты:")
	for {
		fmt.Scanln(first)
		if coinList[*first].Name != "" {
			break
		}
		fmt.Println("Введено некорректное наименование, повторите ввод!1")
	}
	fmt.Print("Введите наименование валюты, в которую необходимо конвертировать:")
	for {
		fmt.Scanln(second)
		if coinList[*first].Name != "" {
			break
		}
		fmt.Println("Введено некорректное наименование, повторите ввод!")
	}
	fmt.Print("Введите количество конвертируемой валюты:")
	for {
		_, err := fmt.Scanf("%f", amount)
		if err == nil {
			break
		}
		fmt.Println("Введено некорректное значение, повторите ввод!2")
	}
}

func calculateRez(first, second *string, amount *float64) float64 {
	var rez float64 = 0

	rez = (coinList[*first].Value * (*amount) / float64(coinList[*first].Nominal)) / (coinList[*second].Value / float64(coinList[*second].Nominal))

	return rez
}

func Project3() {
	var first, second string
	var amount float64
	err := getMoney()

	if err != nil {
		fmt.Println(err)
		return
	}
	readConsole(&first, &second, &amount)
	rez := calculateRez(&first, &second, &amount)
	fmt.Println("При конвертации валюты ", first, " вы получите ", rez, " валюты", second)
}
