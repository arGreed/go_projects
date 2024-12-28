package swTaskManager

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Получить информацию о себе
// @Description Возвращает информацию о текущем пользователе.
// @Tags user
// @Produce json
// @Success 200 {object} Claims
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /showMe [get]
func showMe(c *gin.Context) {
	var claims Claims
	a, _ := c.Get("userId")
	claims.UserId = a.(int64)
	a, _ = c.Get("userPass")
	claims.UserPassword = a.(string)
	a, _ = c.Get("userName")
	claims.UserName = a.(string)
	a, _ = c.Get("userEmail")
	claims.UserEmail = a.(string)
	a, _ = c.Get("userRole")
	claims.UserRole = a.(string)
	c.JSON(http.StatusOK, gin.H{"you": claims})
}

// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body User true "Данные пользователя"
// @Success 200 {object} User
// @Failure 400 {object} map[string]string
// @Router /register [post]
func register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		err := c.ShouldBindJSON(&user)
		if err != nil || !validate(&user) {
			log.Println("Получены данные в ненадлежащем формате.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "corrupted json passed"})
			return
		}
		result := db.Create(&user)
		if result.Error != nil {
			log.Println(result.Error)
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}
		c.JSON(http.StatusOK, gin.H{"User added": user})
	}
}

// @Summary Логин пользователя
// @Description Аутентифицирует пользователя и возвращает JWT-токен.
// @Tags auth
// @Accept json
// @Produce json
// @Param login body Login true "Данные для входа"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /login [post]
func login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var login Login
		var user User
		err := c.ShouldBindJSON(&login)
		if err != nil || !validate(&login) {
			log.Println("Получены данные в ненадлежащем формате.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "corrupted json passed"})
			return
		}
		result := db.Where("email = ? and password = ?", login.Email, login.Password).First(&user)
		if result.Error != nil {
			log.Println(result.Error)
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}
		token, err := generateJWT(&user)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

// @Summary Проверка соединения
// @Description Проверяет, что сервер работает.
// @Tags utils
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "connected"})
}

// @Summary Добавить задачу
// @Description Создает новую задачу.
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body Task true "Данные задачи"
// @Success 200 {object} Task
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /task/add [post]
func addTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var task Task
		err := c.ShouldBindJSON(&task)
		if err != nil || !validate(&task) {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		a, _ := c.Get("userId")
		task.AuthorId = a.(int64)
		result := db.Create(&task)
		if result.Error != nil {
			log.Println(result.Error)
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}

		c.JSON(http.StatusOK, gin.H{"task": task})
	}
}

// @Summary Получить список всех задач
// @Description Возвращает список всех задач (доступно только админу).
// @Tags tasks
// @Produce json
// @Success 200 {array} Task
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /task/all [get]
func allTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("userRole")
		if role.(string) != string("admin") {
			log.Println("Попытка вызова защищённого маршрута.")
			c.JSON(http.StatusLocked, gin.H{"error": "У вас недостаточно прав!"})
			return
		}
		var tasks []Task
		result := db.Find(&tasks)
		if result.Error != nil {
			log.Println(result.Error)
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}
		c.JSON(http.StatusOK, gin.H{"tasks": tasks})
	}
}

// @Summary Получить список всех пользователей
// @Description Возвращает список всех пользователей (доступно только админу).
// @Tags user
// @Produce json
// @Success 200 {array} User
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /user/all [get]
func allUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("userRole")
		if role.(string) != string("admin") {
			log.Println("Попытка вызова защищённого маршрута.")
			c.JSON(http.StatusLocked, gin.H{"error": "У вас недостаточно прав!"})
			return
		}
		var users []User
		result := db.Find(&users)
		if result.Error != nil {
			log.Println(result.Error)
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}
		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

// @Summary Удалить пользователя
// @Description Удаляет пользователя по его ID (доступно только админу).
// @Tags user
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /user/delete/{id} [delete]
func delUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("userRole")
		if role.(string) != string("admin") {
			log.Println("Попытка вызова защищённого маршрута.")
			c.JSON(http.StatusLocked, gin.H{"error": "У вас недостаточно прав!"})
			return
		}
		param := c.Param("id")
		if param == "" {
			log.Println("Некорректный вызов")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Отсутствует идентификатор удаляемого пользователя"})
			return
		}
		usrId, err := strconv.Atoi(param)
		if err != nil {
			log.Println("Передан параметр некорректного типа")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Передан параметр некорректного типа"})
			return
		}
		var user User
		result := db.Find(&user, usrId)
		if result.Error != nil {
			log.Println(result.Error)
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}
		result = db.Delete(&User{}, usrId)
		if result.Error != nil {
			log.Println(result.Error)
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}
		c.JSON(http.StatusGone, gin.H{"deleted user": user})
	}
}

// @Summary Удалить задачу
// @Description Удаляет задачу по её ID.
// @Tags tasks
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /task/delete/{id} [delete]
func delTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("id")
		if param == "" {
			log.Println("Некорректный вызов")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Отсутствует идентификатор удаляемого пользователя"})
			return
		}
		taskId, err := strconv.Atoi(param)
		if err != nil {
			log.Println("Передан параметр некорректного типа")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Передан параметр некорректного типа"})
			return
		}
		curUsrId, _ := c.Get("userId")
		curUsrRole, _ := c.Get("userRole")
		var task Task
		result := db.First(&task, taskId)
		if result.Error != nil {
			log.Println(result.Error)
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}

		if !(curUsrId.(int64) == task.AuthorId || curUsrRole.(string) == "admin") {
			log.Println("Попытка удаления чужой задачи.")
			c.JSON(http.StatusLocked, gin.H{"error": "Попытка удаления чужой задачи."})
			return
		}

		result = db.Delete(&Task{}, taskId)
		if result.Error != nil {
			log.Println(result.Error)
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}
		c.JSON(http.StatusGone, gin.H{"deleted task": task})
	}
}

// @Summary Обновить задачу
// @Description Обновляет задачу по её ID.
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param task body Task true "Новые данные задачи"
// @Success 200 {object} Task
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /task/update/{id} [put]
func updateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("id")
		if param == "" {
			log.Println("Некорректный вызов")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Отсутствует идентификатор удаляемого пользователя"})
			return
		}
		taskId, err := strconv.Atoi(param)
		if err != nil {
			log.Println("Передан параметр некорректного типа")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Передан параметр некорректного типа"})
			return
		}
		curUsrId, _ := c.Get("userId")
		curUsrRole, _ := c.Get("userRole")
		var task Task
		result := db.First(&task, taskId)
		if result.Error != nil {
			log.Println(result.Error)
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}

		if !(curUsrId.(int64) == task.AuthorId || curUsrRole.(string) == "admin") {
			log.Println("Попытка изменения чужой задачи.")
			c.JSON(http.StatusLocked, gin.H{"error": "Попытка изменения чужой задачи."})
			return
		}

		var newTask Task

		err = c.ShouldBindJSON(&newTask)

		if err != nil || !validate(&newTask) {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		newTask.Id = task.Id
		newTask.AuthorId = task.AuthorId
		db.Model(&task).Updates(newTask)
		c.JSON(http.StatusOK, gin.H{"new task": newTask})
	}
}
