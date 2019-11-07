package models

type RawData struct {
	UserId          string
	Email           string
	Password        string
	Bill            float64
	ItemId          string
	ItemName        string
	InventoryUserId string
	InventoryItemId string
	InventoryCount  int64
	SalesId         string
	SalesUserId     string
	ItemSaleSaleId  string
	ItemSaleItemId  string
	ItemSaleCount   int64
	ItemSalePrice   float64
}
