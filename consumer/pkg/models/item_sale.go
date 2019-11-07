package models

import "github.com/satori/go.uuid"

type ItemSale struct {
	SellerID uuid.UUID `gorm:"-"`
	SaleID   uuid.UUID `gorm:"primary_key"`
	ItemID   uuid.UUID `gorm:"primary_key"`
	Name     string    `gorm:"-"`
	Count    int64
	Price    float64
}
