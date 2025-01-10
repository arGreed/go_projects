package firstChat

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	ping     string = "/ping"
	register string = "/register"
)

func getPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "connected"})
}

func getRegister(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
