package models

import "github.com/satori/go.uuid"

type Sales struct {
	ID     uuid.UUID `gorm:"primary_key"`
	UserID uuid.UUID `gorm:"column:user_id"`
}
