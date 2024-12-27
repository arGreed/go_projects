package swTaskManager

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dsn       string = "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable"
	logFile   string = "project23/test.log"
	jwtSecret string = "EVDFOIGHNWECNOGVI[iodnhbplsefnqweoijnmgherpsdopiawmf['wsedm,g[]]]"
)

func generateJWT(user *User) (string, error) {
	claims := &Claims{
		UserId:       user.Id,
		UserName:     user.Name,
		UserEmail:    user.Email,
		UserPassword: user.Password,
		UserRole:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		log.Println("Ошибка генерации токена!")
		return "", errors.New("ошибка генерации токена")
	}

	return tokenString, err

}

func prepareDB(db *gorm.DB) error {

	err := db.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&Task{})
	if err != nil {
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

func logInit() (*os.File, error) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return nil, err
	}
	log.SetOutput(file)
	return file, err
}

func logMiddleware(c *gin.Context) {
	log.Println(c.Request.Method, c.Request.URL.Path)
}

type validator interface {
	isValid() bool
}

func validate(v validator) bool {
	return v.isValid()
}

func (user User) isValid() bool {
	if user.Email == "" || !strings.Contains(user.Email, "@") || user.Name == "" || user.Password == "" {
		return false
	}
	return true
}

func (login Login) isValid() bool {
	if login.Email == "" || !strings.Contains(login.Email, "@") || login.Password == "" {
		return false
	}
	return true
}

func (task Task) isValid() bool {
	if task.Description == "" || task.Name == "" {
		return false
	}
	return true
}

func authMiddleware(c *gin.Context) {
	AuthToken := c.GetHeader("Authentication")
	if AuthToken == "" {
		log.Println("Попытка перейти по защищённому маршруту без авторизации.")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Сперва необходимо авторизироваться"})
		return
	}
	tokenString := AuthToken[len("Bearer "):]

	claims := &Claims{}
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

	c.Set("userId", claims.UserId)
	c.Set("userPass", claims.UserPassword)
	c.Set("userName", claims.UserName)
	c.Set("userEmail", claims.UserEmail)
	c.Set("userRole", claims.UserRole)

	c.Next()
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
	// Проверка коннекции.
	router.GET("/ping", logMiddleware, ping)
	// Регистрация.
	router.POST("/register", logMiddleware, register(db))
	// Логин.
	router.POST("/login", logMiddleware, login(db))
	// Получить информацию о себе.
	router.GET("/showMe", logMiddleware, authMiddleware, showMe)
	// Добавить задачу в базу.
	router.POST("/task/add", logMiddleware, authMiddleware, addTask(db))
	// Админ список задач.
	router.GET("/task/all", logMiddleware, authMiddleware, allTask(db))
	// Админ список пользователей.
	router.GET("/user/all", logMiddleware, authMiddleware, allUser(db))
	// Админ удаление пользователя.
	router.DELETE("/user/delete/:id", logMiddleware, authMiddleware, delUser(db))
	router.DELETE("/task/delete/:id", logMiddleware, authMiddleware, delTask(db))
	router.PUT("/task/update/:id", logMiddleware, authMiddleware, updateTask(db))

	router.Run("localhost:8081")
}
