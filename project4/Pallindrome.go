package project4

import (
	"fmt"
	"strings"
)

var (
	stringsStorage []string
)

func read() {
	var s string
	fmt.Println("Введите проверяемые строки")
	for {
		fmt.Scanln(&s)
		if s != "" {
			s = strings.ReplaceAll(s, " ", "")
			s = strings.ToLower(s)
			stringsStorage = append(stringsStorage, s)
		} else {
			break
		}
		s = ""
	}
}

func isPalindrome(s string) bool {
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] {
			return false
		}
	}

	return true
}

func check() int {
	var rez int = 0
	for _, i := range stringsStorage {
		if isPalindrome(i) {
			rez++
		}
	}

	return rez
}

func Project4() {
	read()
	var rez int = check()
	fmt.Println("Количество введённых палиндромов: ", rez)
}
