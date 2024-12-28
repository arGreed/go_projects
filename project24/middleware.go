package weather

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = "gwepgjspovwjr9ujhbdfinvadwkm0-[gjgreiobne[kw,[ofneroinermfol]]]"

func generateJWT(user *User) (string, error) {
	claims := Claims{
		UserId:       user.Id,
		UserName:     user.Name,
		UserEmail:    user.Email,
		UserPassword: user.Password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func logMiddleware(c *gin.Context) {
	log.Println(c.Request.Method, c.Request.URL.Path)
}

func authMiddleware(c *gin.Context) {
	AuthToken := c.GetHeader("Authentication")
	if AuthToken == "" {
		log.Println("Попытка перейти по защищённому маршруту без авторизации.")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Сперва необходимо авторизироваться"})
		return
	}
	tokenString := AuthToken[len("Bearer "):]

	var claims Claims

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "неверная подпись токена"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный токен"})
		}
		c.Abort()
		return
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный токен"})
		c.Abort()
		return
	}
	c.Next()
}
