package curseController

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Структура пользователя.
type User struct {
	Id       int    `gorm:"primaryKey" json:"-"`
	Age      int    `gorm:"unique; not null" json:"Age"`
	Name     string `gorm:"unique; not null" json:"Name"`
	Email    string `gorm:"unique; not null" json:"Email"`
	Password string `gorm:"not null" json:"Password"`
	Role     string `gorm:"default:'user'" json:"Role"`
}

// Структура курса.
type Course struct {
	Id       int       `gorm:"primaryKey" json:"-"`
	Capacity int       `gorm:"not null" json:"Capacity"`
	Name     string    `gorm:"unique; not null" json:"Name"`
	Title    string    `gorm:"not null" json:"Title"`
	startDt  time.Time `json:"startDt"`
	endDt    time.Time `json:"endDt"`
}

// Структура записи на курс.
type Enrollment struct {
	Id       int       `gorm:"primaryKey" json:"Id"`
	UserId   uint      `json:"UserId"`
	CourseId uint      `json:"CourseId"`
	Date     time.Time `gorm:"default:current_timestamp" json:"Date"`
}

// Структура авторизированного пользователя.
type Claims struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Структура для авторизации пользователя
type LoginParams struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}
