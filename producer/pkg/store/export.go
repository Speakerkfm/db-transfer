package store

import (
	"db-transfer/producer/pkg/models"
)

func (st *Store) GetData() ([]models.RawData, error) {
	rows, err := st.db.Query("SELECT * FROM shop;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
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
	)

	var rawData []models.RawData

	i := 0
	for rows.Next() {
		err = rows.Scan(&UserId, &Email, &Password, &Bill, &ItemId, &ItemName, &InventoryUserId, &InventoryItemId, &InventoryCount, &SalesId, &SalesUserId, &ItemSaleSaleId, &ItemSaleItemId, &ItemSaleCount, &ItemSalePrice)
		if err != nil {
			return nil, err
		}

		rawData = append(rawData, models.RawData{
			UserId:          UserId,
			Email:           Email,
			Password:        Password,
			Bill:            Bill,
			ItemId:          ItemId,
			ItemName:        ItemName,
			InventoryUserId: InventoryUserId,
			InventoryItemId: InventoryItemId,
			InventoryCount:  InventoryCount,
			SalesId:         SalesId,
			SalesUserId:     SalesUserId,
			ItemSaleSaleId:  ItemSaleSaleId,
			ItemSaleItemId:  ItemSaleItemId,
			ItemSaleCount:   ItemSaleCount,
			ItemSalePrice:   ItemSalePrice,
		})

		i++
	}

	return rawData, nil
}
