package firstChat

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	logFile   string = "project25/test.log"
	secretJWT string = "weoignSPOGJM223490GNOFEGN 209IU23-Q0GNIFN  1Q02IEJH1209R0-3GFJMWEGNM0PQ23RJM"
	dsn       string = "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable"
)

func logInit() (*os.File, error) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		return nil, err
	}
	log.SetOutput(file)
	return file, nil
}

func dbInit() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func generateJWT(user *User) string {
	var tokenString string
	claims := Claims{
		userId:   user.Id,
		userRole: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return tokenString
}

func WebChat() {
	// Инициализация лог файла.
	file, err := logInit()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Отложенное закрытие файла логов.
	defer file.Close()
	// Инициализация коннекции к БД.
	// db, err := dbInit()
	if err != nil {
		log.Println(err)
		return
	}

	router := gin.Default()
	// Проверка соединения.
	router.GET(ping, logMiddleware, getPing)
	// Регистрация.
	router.POST(register, logMiddleware, getRegister(db))

	router.Run("localhost:8081")
}
