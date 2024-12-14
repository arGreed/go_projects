package project5

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	inpTxt string
)

func read() error {
	scanner := bufio.NewScanner(os.Stdin)
	var buf string
	fmt.Println("Введите текст, в котором необходимо посчитать количество слов")
	for scanner.Scan() {
		buf = scanner.Text()
		if buf != "" {
			inpTxt += buf
		} else {
			break
		}
	}
	return nil
}

func countWords() int {
	words := strings.Fields(inpTxt)
	return len(words)
}

func Project5() {
	read()
	fmt.Println("Количество слов в ведённом тексте равно ", countWords())
}
