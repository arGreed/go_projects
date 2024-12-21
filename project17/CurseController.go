package curseController

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

func usrIsValid(usr *User) bool {
	if usr.Id != 0 || usr.Age < 1 || usr.Age > 100 || !strings.Contains(usr.Email, "@") || usr.Email == "" || usr.Name == "" || usr.Password == "" || usr.Role == "" {
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
		if !usrIsValid(&user) {
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

func CurseController() {

	db, err := dbInit()
	if err != nil {
		return
	}
	router := gin.Default()

	router.POST("/user", addUser(db))

	router.Run("localhost:8081")
}
