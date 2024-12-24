package feedBack

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var logFile string = "project21/test.log"

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"info": "server connected"})
}

func logMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf(c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}

func logInit() (*os.File, error) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		return nil, err
	}
	log.SetOutput(file)
	return file, nil
}

func submit(c *gin.Context) {
	var feedback UsrFeedBack
	if err := c.ShouldBind(&feedback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вывод данных в консоль (или сохранение в базу данных)
	fmt.Printf("Received Feedback: %+v\n", feedback)

	// Ответ пользователю
	c.JSON(http.StatusOK, gin.H{"message": "Feedback received successfully!"})
}

func showForm(c *gin.Context) {
	c.HTML(http.StatusOK, "feedback.html", nil)
}

func FeedBack() {
	logs, err := logInit()
	if err != nil {
		fmt.Println("Произошла ошибка при инициализации файла логов", err)
		return
	}
	defer logs.Close()

	router := gin.Default()
	router.LoadHTMLFiles("project21/feedback.html")
	router.GET("/ping", logMiddleware(), ping)
	router.GET("/feedback", showForm)
	router.POST("/submit", submit)

	router.Run("localhost:8081")
}
