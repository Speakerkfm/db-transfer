package services

import (
	"db-transfer/consumer/pkg/models"
	"db-transfer/consumer/pkg/store"
	"encoding/json"
)

type Transfer struct {
	st   *store.Store
}

func NewTransfer(st *store.Store) *Transfer {
	return &Transfer{
		st:   st,
	}
}

func (t *Transfer) Handle(bytes []byte) error {
	var data []models.RawData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	return t.st.ImportData(data)
}
