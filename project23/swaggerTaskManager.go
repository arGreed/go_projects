package swTaskManager

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dsn     string = "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable"
	logFile        = "project23/test.log"
)

func dbInit() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func logInit() (*os.File, error) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return nil, err
	}
	log.SetOutput(file)
	return file, err
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "connected"})
}

func logMiddleware(c *gin.Context) {
	log.Println(c.Request.Method, c.Request.URL.Path)
}

func register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func SwTask() {
	logPtr, err := logInit()
	if err != nil {
		fmt.Println("error log init")
		return
	}
	defer logPtr.Close()
	db, err := dbInit()
	if err != nil {
		fmt.Println("error db init")
		return
	}

	router := gin.Default()

	router.GET("/ping", logMiddleware, ping)
	router.POST("/register", logMiddleware, register(db))

	router.Run("localhost:8081")
}
