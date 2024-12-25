package grcoder

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/skip2/go-qrcode"
)

func grGenerate(s *string) {
	qrCode, err := qrcode.New(*s, qrcode.Medium)
	if err != nil {
		log.Fatal(err)
	}
	err = qrCode.WriteFile(256, "qrcode.png")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("QR-код успешно создан и сохранен в файл qrcode.png")
}

func GrCoder() {
	scanner := bufio.NewScanner(os.Stdin)
	var s string
	fmt.Print("Введите сообщение, которое хотите закодировать:")

	if scanner.Scan() {
		s = scanner.Text()
	}
	if s == "" {
		return
	}
	grGenerate(&s)
}
