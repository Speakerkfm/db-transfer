package store

import (
	"db-transfer/consumer/pkg/models"
	"github.com/satori/go.uuid"
)

func (st *Store) ImportData(data []models.RawData) error {
	for _, row := range data {
		user := models.User{
			ID: uuid.FromStringOrNil(row.UserId),
			Email: row.Email,
			Password: row.Password,
			Bill: row.Bill,
		}
		if err := st.gorm.FirstOrCreate(&user).Error; err != nil {
			return err
		}

		item := models.Item{
			ID: uuid.FromStringOrNil(row.ItemId),
			Name: row.ItemName,
		}
		if err := st.gorm.FirstOrCreate(&item).Error; err != nil {
			return err
		}

		inventory := models.Inventory{
			UserID: uuid.FromStringOrNil(row.InventoryUserId),
			ItemID: uuid.FromStringOrNil(row.InventoryItemId),
			Count: row.InventoryCount,
		}
		if err := st.gorm.FirstOrCreate(&inventory).Error; err != nil {
			return err
		}

		sale := models.Sales{
			ID: uuid.FromStringOrNil(row.SalesId),
			UserID: uuid.FromStringOrNil(row.SalesUserId),
		}
		if err := st.gorm.FirstOrCreate(&sale).Error; err != nil {
			return err
		}

		itemSale := models.ItemSale{
			SaleID: uuid.FromStringOrNil(row.ItemSaleSaleId),
			ItemID: uuid.FromStringOrNil(row.ItemSaleItemId),
			Count: row.ItemSaleCount,
			Price: row.ItemSalePrice,
		}
		if err := st.gorm.FirstOrCreate(&itemSale).Error; err != nil {
			return err
		}
	}

	return nil
}
