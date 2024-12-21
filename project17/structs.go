package curseController

import (
	"time"
)

type User struct {
	Id       int    `gorm:"primaryKey" json:"-"`
	Age      int    `gorm:"unique; not null" json:"Age"`
	Name     string `gorm:"unique; not null" json:"Name"`
	Email    string `gorm:"unique; not null" json:"Email"`
	Password string `gorm:"not null" json:"Password"`
	Role     string `gorm:"default:'user'" json:"Role"`
}

type Course struct {
	Id       int       `gorm:"primaryKey" json:"Id"`
	Capacity int       `gorm:"not null" json:"Capacity"`
	Name     string    `gorm:"unique; not null" json:"Name"`
	Title    string    `gorm:"not null" json:"Title"`
	startDt  time.Time `json:"startDt"`
	endDt    time.Time `json:"endDt"`
}

type Enrollment struct {
	Id       int       `gorm:"primaryKey" json:"Id"`
	UserId   uint      `json:"UserId"`
	CourseId uint      `json:"CourseId"`
	Date     time.Time `gorm:"default:current_timestamp" json:"Date"`
}
