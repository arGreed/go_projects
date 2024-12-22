package simpleWeb

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var logName string = "project20/test.log"

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "server ok"})
	c.Request.Body.Close()
}

func logInit() (*os.File, error) {
	file, err := os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	log.SetOutput(file)
	return file, nil
}

func logMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.Request.Method, c.Request.URL.Path)
		c.JSON(http.StatusOK, gin.H{"status": "server ok"})
		c.Request.Body.Close()
	}
}

func WebServer() {
	log, err := logInit()
	if err != nil {
		fmt.Println()
	}
	defer log.Close()

	router := gin.Default()

	router.GET("/ping", logMiddleware(), ping)

	router.Run("localhost:8081")
}
