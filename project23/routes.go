package swTaskManager

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "connected"})
}

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
