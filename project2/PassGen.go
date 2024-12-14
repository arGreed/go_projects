package project2

import (
	"fmt"
	"math/rand"
	"strconv"
)

func read(len *int, mode *string) error {
	var buf string
	var err error
	for {
		fmt.Print("Введите набор символов для генерации пароля (ru, eng):")
		fmt.Scanln(&buf)
		if buf == "ru" || buf == "eng" {
			*mode = buf
			break
		}
		fmt.Println("Выбран некорректный режим, повторите ввод!")
	}
	for {
		fmt.Print("Ведите длину строки генерируемого пароля (1-100 символов):")
		fmt.Scanln(&buf)
		*len, err = strconv.Atoi(buf)
		if err == nil && *len > 0 && *len < 101 {
			break
		}
		fmt.Println("Введено некорректное значение, повторите ввод!")
	}

	return nil
}

var (
	eng string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rus string = "абвгдёежзийклмнопрстуфхцчшщъыьэюяАБВГДЁЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ123456789"
)

func pwdGenerate(pwd, mode *string, ln *int) {
	result := make([]byte, *ln)
	if *mode == "ru" {
		for i := range result {
			result[i] = rus[rand.Intn(len(rus))]
		}
	} else {
		for i := range result {
			result[i] = eng[rand.Intn(len(eng))]
		}
	}
	*pwd = string(result)
}

func Project2() {
	var len int
	var mode, pwd string
	read(&len, &mode)
	pwdGenerate(&pwd, &mode, &len)

	fmt.Println("Сгенерированный пароль: ", pwd)
}
