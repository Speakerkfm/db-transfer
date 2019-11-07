package models

import "github.com/satori/go.uuid"

type Inventory struct {
	UserID uuid.UUID `gorm:"primary_key"`
	ItemID uuid.UUID `gorm:"primary_key"`
	Name   string    `gorm:"-"`
	Count  int64
}
