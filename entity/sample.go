package entity

import "github.com/lib/pq"

type User struct {
	ID      int            `json:"id"`
	Age     int            `json:"age"`
	Name    string         `json:"name"`
	Hobbies pq.StringArray `json:"hobbies" gorm:"type:text[]"`
	Jobs    []Job          `json:"jobs"`
}

type Job struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserId int    `json:"user_id"`
}
