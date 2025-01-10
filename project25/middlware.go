package firstChat

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Логирование маршрутов, по которым ходит пользователь.
func logMiddleware(c *gin.Context) {
	var s string
	s = c.Request.Method
	s += c.Request.URL.Path
	log.Println(s)
}
