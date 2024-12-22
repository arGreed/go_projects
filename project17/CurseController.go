package curseController

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	jwtSecret string = "piertogshderjkweigfhnS]ERHBJK-2W9GJHW0EVINWEROW3HN4G0PW3ENFV"
)

func generateJWT(user *User) (string, error) {
	claims := &Claims{
		UserId: user.Id,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Срок действия токена - 24 часа
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // Время выдачи токена
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))

	return tokenString, err
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "токен отсутствует"})
			c.Abort()
			return
		}

		// Извлекаем токен из заголовка
		tokenString := authHeader[len("Bearer "):]

		// Парсим токен
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

		// Сохраняем данные пользователя в контексте запроса
		c.Set("userId", claims.UserId)
		c.Set("role", claims.Role)

		// Продолжаем выполнение запроса
		c.Next()
	}

}

func prepareDB(db *gorm.DB) error {

	err := db.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&Course{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&Enrollment{})
	if err != nil {
		return err
	}
	return nil
}

func dbInit() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = prepareDB(db)

	return db, nil
}

type validator interface {
	IsValid() bool
}

func validate(v validator) bool {
	return v.IsValid()
}

func (user User) IsValid() bool {
	if user.Id != 0 || user.Age < 1 || user.Age > 100 || !strings.Contains(user.Email, "@") || user.Email == "" || user.Name == "" || user.Password == "" || user.Role == "" || user.Role == "admin" {
		return false
	}
	return true
}

func (course Course) IsValid() bool {
	if course.endDt.IsZero() || course.startDt.IsZero() || course.Capacity == 0 || course.Name == "" || course.Title == "" {
		return false
	}
	return true
}

func (logPrms LoginParams) IsValid() bool {
	if logPrms.Email == "" || !strings.Contains(logPrms.Email, "@") || logPrms.Password == "" {
		return false
	}
	return true
}

func addUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			c.Request.Body.Close()
			return
		}
		if !validate(user) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "получены некорректные данные"})
			c.Request.Body.Close()
			return
		}
		result := db.Create(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			c.Request.Body.Close()
			return
		}
		c.JSON(http.StatusCreated, gin.H{"added user": user})
	}
}

func addCurse(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var course Course
		err := c.ShouldBindJSON(&course)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			c.Request.Body.Close()
			return
		}
		if validate(course) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "получены некорректные данные"})
			c.Request.Body.Close()
			return
		}
		result := db.Create(&course)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			c.Request.Body.Close()
			return
		}
		c.JSON(http.StatusCreated, gin.H{"added course": course})
	}
}

func admAllUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role.(string) != "admin" {
			c.JSON(http.StatusLocked, gin.H{"error": "недостаточно прав."})
			c.Request.Body.Close()
			return
		}
		var users []User

		result := db.Find(&users)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			c.Request.Body.Close()
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

func admDeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role.(string) != "admin" {
			c.JSON(http.StatusLocked, gin.H{"error": "недостаточно прав."})
			c.Request.Body.Close()
			return
		}

		param := c.Param("id")
		if param == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Параметр не найден"})
			c.Request.Body.Close()
			return
		}
		var UserId, _ = strconv.Atoi(param)

		var user User

		result := db.Delete(&user, UserId)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			c.Request.Body.Close()
			return
		}

		c.JSON(http.StatusOK, gin.H{"Info": "Deleted"})
	}
}

func admAllCourses(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role.(string) != "admin" {
			c.JSON(http.StatusLocked, gin.H{"error": "недостаточно прав."})
			c.Request.Body.Close()
			return
		}
		var courses []Course

		result := db.Find(&courses)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			c.Request.Body.Close()
			return
		}

		c.JSON(http.StatusOK, gin.H{"courses": courses})
	}
}

func admDeleteCourse(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role.(string) != "admin" {
			c.JSON(http.StatusLocked, gin.H{"error": "недостаточно прав."})
			c.Request.Body.Close()
			return
		}

		param := c.Param("id")
		if param == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Параметр не найден"})
			c.Request.Body.Close()
			return
		}
		var CourseId, _ = strconv.Atoi(param)

		var course Course

		result := db.Delete(&course, CourseId)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			c.Request.Body.Close()
			return
		}

		c.JSON(http.StatusOK, gin.H{"Info": "Deleted"})
	}
}

func login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginParams LoginParams
		var user User

		// Привязываем JSON к структуре loginParams
		err := c.ShouldBindJSON(&loginParams)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Валидируем данные
		if !validate(loginParams) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "получены некорректные данные"})
			return
		}

		// Ищем пользователя в базе данных
		result := db.Where("email = ? AND password = ?", loginParams.Email, loginParams.Password).First(&user)

		// Проверяем, была ли ошибка
		if result.Error != nil {
			// Если пользователь не найден, возвращаем 401
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный email или пароль"})
				return
			}

			// Если другая ошибка, возвращаем 500
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		// Генерируем JWT-токен
		token, err := generateJWT(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Возвращаем токен
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func userEnroll(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("id")
		if param == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Параметр не найден"})
			c.Request.Body.Close()
			return
		}
		var courseId, _ = strconv.Atoi(param)
		var course Course
		result := db.First(&course, courseId)

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			c.Request.Body.Close()
			return
		}
		var cnt int64 = 0
		result = db.Model(&Enrollment{}).Where("course_id = ?", courseId).Count(&cnt)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			c.Request.Body.Close()
			return
		}
		if cnt >= int64(course.Capacity) {
			c.JSON(http.StatusConflict, gin.H{"error": "Вместимость курса превышена."})
			c.Request.Body.Close()
			return
		}

		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
			c.Request.Body.Close()
			return
		}
		enrollment := Enrollment{
			UserId:   userId.(uint),
			CourseId: uint(cnt),
			Date:     time.Now(),
		}
		result = db.Create(&enrollment)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при записи на курс"})
			c.Request.Body.Close()
			return
		}

		// Возвращаем успешный ответ
		c.JSON(http.StatusCreated, gin.H{"message": "Пользователь успешно записан на курс"})
	}
}

func userUpdate(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("id")
		if param == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Параметр не найден"})
			c.Request.Body.Close()
			return
		}
		var user User
		userId, _ := c.Get("userId")
		result := db.First(&user, userId.(int))

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			c.Request.Body.Close()
			return
		}

		var newUser User

		err := c.ShouldBindJSON(&newUser)

		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err})
			c.Request.Body.Close()
			return
		}

		if !validate(newUser) {
			c.JSON(http.StatusConflict, gin.H{"error": "Некорректный ввод."})
			c.Request.Body.Close()
			return
		}

		user = newUser

		user.Id = userId.(int)

		db.Save(&user)
		// Возвращаем успешный ответ
		c.JSON(http.StatusCreated, gin.H{"message": "Параметры успешно изменены!"})
	}
}

func CurseController() {

	db, err := dbInit()
	if err != nil {
		return
	}
	router := gin.Default()

	router.GET("/login", login(db))

	// Операции администратора.
	router.GET("/users", authMiddleware(), admAllUsers(db))
	router.DELETE("/users/delete/:id", authMiddleware(), admDeleteUser(db))
	router.GET("/courses", authMiddleware(), admAllCourses(db))
	router.DELETE("/courses/delete/id", authMiddleware(), admDeleteCourse(db))

	// Операции с пользователями.
	// Доступно всем.
	router.POST("/user", addUser(db))
	router.POST("/user/enroll/:id", authMiddleware(), userEnroll(db))
	router.POST("/user/update", authMiddleware(), userUpdate(db))
	// Операции с курсами.
	// Только с авторизацией.
	router.POST("/course", authMiddleware(), addCurse(db))

	router.Run("localhost:8081")
}
