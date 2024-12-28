package weather

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "connected"})
}

func register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		err := c.ShouldBindJSON(&user)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		if !validate(&user) {
			log.Println("Получены некорректные данные.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Получены некорректные данные."})
			return
		}
		result := db.Create(&user)
		if result.Error != nil {
			log.Println(result.Error)
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var logParams LogParams
		var user User
		err := c.ShouldBindJSON(&logParams)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		result := db.Where("email = ? and password = ?", logParams.Email, logParams.Password).First(&user)
		if result.Error != nil {
			log.Println(result.Error)
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}
		tokenString, err := generateJWT(&user)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}

func getWeather(c *gin.Context) {
	city := c.Param("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "параметр был передан"})
		return
	}
}
