package firstChat

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Message struct {
	Sender  string `json:"Sender"`
	Context string `json:"Context"`
}

type User struct {
	Id       int64  `json:"-" gorm:"primary key"`
	Name     string `json:"Name" gorm:"type:varchar(100); not null"`
	Email    string `json:"Email" gorm:"type:varchar(100); not null"`
	Password string `json:"Password" gorm:"type:varchar(100); not null"`
	Role     string `json:"-" gorm:"type:varchar(50);default: 'user'"`
}

type Claims struct {
	userId   int64
	userRole string
	jwt.RegisteredClaims
}

type Login struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type validator interface {
	isValid() bool
}

func validate(v validator) bool {
	return v.isValid()
}

func (user User) isValid() bool {
	if !strings.Contains(user.Email, "@") || user.Id != 0 || user.Name == "" || len(user.Password) <= 5 || user.Role != "" {
		return false
	}
	return true
}

func (login Login) isValid() bool {
	if !strings.Contains(login.Email, "@") || len(login.Password) <= 5 {
		return false
	}
	return true
}
