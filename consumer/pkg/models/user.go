package models

import "github.com/satori/go.uuid"

type User struct {
	ID       uuid.UUID `gorm:"primary_key"`
	Email    string
	Password string
	Bill     float64
}

func (User) TableName() string {
	return "users"
}
