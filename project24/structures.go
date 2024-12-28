package weather

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type validator interface {
	isValid() bool
}

func validate(v validator) bool {
	return v.isValid()
}

func (user User) isValid() bool {
	if !strings.Contains(user.Email, "@") || user.Password == "" || user.Name == "" {
		return false
	}
	return true
}

type User struct {
	Id       int64  `json:"-" gorm:"primaryKey"`
	Name     string `json:"Name" gorm:"unique; not null"`
	Email    string `json:"Email", gorm:"unique; not null"`
	Password string `json:"Password", gorm:"not null"`
}

type Claims struct {
	UserId               int64  `json:"UserId"`
	UserName             string `json:"UserName"`
	UserEmail            string `json:"UserEmail"`
	UserPassword         string `json:"UserPassword"`
	jwt.RegisteredClaims `json:"-"`
}

type LogParams struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

package main

import (
    "fmt"
    "log"

    "github.com/go-resty/resty/v2"
)
