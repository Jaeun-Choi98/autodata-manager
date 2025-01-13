package entity

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID      int            `json:"id"`
	Age     int            `json:"age"`
	Name    string         `json:"name"`
	Hobbies pq.StringArray `json:"hobbies" gorm:"type:text[]"`
	Jobs    []Job          `json:"jobs" gorm:"foreignKey:UserId;references:ID"`
}

// has many
type Job struct {
	gorm.Model
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserId int    `json:"user_id"`
}
