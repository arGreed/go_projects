package swTaskManager

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Id       int64  `json:"-" gorm:"primaryKey"`
	Email    string `json:"Email" gorm:"not null; unique"`
	Password string `json:"Password" gorm:"not null"`
	Name     string `json:"Name", gorm:"not null"`
	Role     string `json:"Role", gorm:"default:'user'"`
}

type Login struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type Claims struct {
	UserId               int64  `json:"user_id"`
	UserName             string `json:"user_name"`
	UserEmail            string `json:"user_email"`
	UserPassword         string `json:"user_password"`
	UserRole             string `json:"user_role"`
	jwt.RegisteredClaims `swaggerignore:"true"`
}

type Task struct {
	Id          int64  `json:"-" gorm:"primaryKey"`
	AuthorId    int64  `json:"-", gorm:"not null"`
	Name        string `json:"Name" gorm:"not null"`
	Description string `json:"Description" gorm:"not null"`
}
