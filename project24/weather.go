package weather

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var logPath string = "project24/test.log"
var dsn string = "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable"

func logInit() (*os.File, error) {
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)

	return file, nil
}

func prepareDB(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func dbInit() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = prepareDB(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Weather() {

	file, err := logInit()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	db, err := dbInit()

	router := gin.Default()

	// Проверка состояния сервера.
	router.GET("/ping", logMiddleware, ping)
	router.POST("/register", logMiddleware, register(db))
	router.GET("/weather/:city", logMiddleware, authMiddleware, getWeather)

	router.Run("localhost:8081")
}
