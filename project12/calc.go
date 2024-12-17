package calculator

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var logFile string = "project12/test.log"
var psblOper string = "*/+-"

// Текст
var (
	inpFPrm string = "Введите первый параметр:"
	inpSPrm string = "Введите второй параметр:"
	inpOper string = "Введите выполняемую операцию:"
	wInp    string = "Некорректный ввод, повторите попытку!"
)

var (
	errNullDiv = errors.New("деление на 0")
	errStrange = errors.New("неожиданная операция")
)

func logInit() (*os.File, error) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	log.SetOutput(file)
	return file, nil
}

func paramReader() (float64, float64, string) {
	scanner := bufio.NewScanner(os.Stdin)
	var first, second float64
	var operation string
	var buf string
	var err error
	for {
		fmt.Println(inpFPrm)
		if scanner.Scan() {
			buf = scanner.Text()
		}
		first, err = strconv.ParseFloat(buf, 64)
		if err == nil {
			break
		}
		fmt.Println(wInp)
		log.Println(wInp)
	}
	for {
		fmt.Println(inpSPrm)
		if scanner.Scan() {
			buf = scanner.Text()
		}
		second, err = strconv.ParseFloat(buf, 64)
		if err == nil {
			break
		}
		fmt.Println(wInp)
		log.Println(wInp)
	}
	for {
		fmt.Println(inpOper)
		if scanner.Scan() {
			buf = scanner.Text()
		}
		if len(buf) == 1 && strings.Contains(psblOper, buf) {
			operation = buf
			break
		}
		fmt.Println(wInp)
		log.Println(wInp)
	}
	return first, second, operation
}

func calculate(first *float64, second *float64, operand *string) (float64, error) {
	switch *operand {
	case "/":
		if *second == float64(0) {
			return 0, errNullDiv
		} else {
			return *first / (*second), nil
		}
	case "*":
		return *first * (*second), nil
	case "+":
		return *first + *second, nil
	case "-":
		return *first - *second, nil
	default:
		return 0, errStrange
	}
}

func Calculator() {
	var rez float64
	file, err := logInit()
	if err != nil {
		fmt.Println(err)
		log.Println(err)
	}
	defer file.Close()
	first, second, operand := paramReader()
	rez, err = calculate(&first, &second, &operand)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
	}
	fmt.Println("Результат вычислений:", rez)
}
