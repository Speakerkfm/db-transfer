package service

import (
	"db-transfer/producer/pkg/models"
	"db-transfer/producer/pkg/store"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)


type SessionService struct {
	st *store.Store
	client *TcpClient
}

type Session struct {
	ClientID string
	Key      string
}

func NewSessionService(client *TcpClient, st *store.Store) *SessionService {
	return &SessionService{
		client: client,
		st: st,
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (ss *SessionService) CreateSession(clientID string) (*Session, error) {
	publicKey, err := ss.client.RequestPublicKey()
	if err != nil {
		return nil, err
	}

	symmetricKey := randStringRunes(128)

	requestData := make(map[string]interface{}, 1)
	requestData["symmetric-key"] = symmetricKey

	fmt.Println(symmetricKey)
	data, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	err = ss.client.SendRSARequest(clientID, publicKey, createSession, data)

	return &Session{
		ClientID: clientID,
		Key: symmetricKey,
	}, err
}

func (ss *SessionService) GetEncryptedDataAES(clientID, symmetricKey, msgType string, data []byte) ([]byte, error) {
	headers := make(map[string]interface{}, 2)
	headers["X-MESSAGE-TYPE"] = msgType
	headers["X-CLIENT-ID"] = clientID

	encryptedData := encrypt(data, symmetricKey)

	msg := &models.Message{
		Headers: headers,
		Data:    encryptedData,
	}

	reqBody, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return reqBody, nil
}
