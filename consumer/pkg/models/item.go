package models

import "github.com/satori/go.uuid"

type Item struct {
	ID   uuid.UUID `gorm:"primary_key"`
	Name string
}

func (Item) TableName() string {
	return "items"
}
