package models

type Message struct {
	Headers map[string]interface{}
	Data    []byte `json:"data"`
}
