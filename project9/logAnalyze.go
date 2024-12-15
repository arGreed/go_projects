package logAnalyze

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

var logName = "project9/test.log"

const er string = "error"

func openOrCreate() (*os.File, error) {
	log, err := os.OpenFile(logName, os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}
	return log, err
}

func ranFillLog() {
	if rand.Intn(100)%2 == 1 {
		log.Printf("error rand is odd")
	} else {
		log.Printf("rand is even")
	}
}

func analyze() (int, int, int, error) {
	var s string
	var errCnt, okCnt, cnt int = 0, 0, 0
	file, err := os.Open(logName)
	if err != nil {
		return 0, 0, 0, err
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		s = scanner.Text()
		if strings.Contains(s, er) {
			errCnt++
		} else {
			okCnt++
		}
		cnt++
	}
	return errCnt, okCnt, cnt, nil
}

func LogAnalyze() {
	file, err := openOrCreate()
	var n int = rand.Intn(1000)

	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}

	for i := 0; i < n; i++ {
		ranFillLog()
	}
	errors, ok, strings, _ := analyze()
	fmt.Println("К-во ошибок в лог файле:", errors, " К-во строк без ошибок в лог файле:", ok, " Всего строк в лог файле:", strings)
}
