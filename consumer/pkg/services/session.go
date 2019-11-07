package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"db-transfer/consumer/pkg/models"
	"db-transfer/consumer/pkg/store"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"hash"
	"time"
)

type SessionService struct {
	st *store.Store
	hash hash.Hash
}

type Session struct {
	ClientID string
	Key      string
}

func NewSessionService(st *store.Store) *SessionService {
	return &SessionService{
		st: st,
		hash: sha256.New(),
	}
}

func (ss *SessionService) CreateSession(msg *models.Message, privateKey *rsa.PrivateKey) error {
	clientID, ok := msg.Headers["X-CLIENT-ID"]
	if !ok {
		return errors.New("Unknown client")
	}

	label := []byte("")
	plainText, err := rsa.DecryptOAEP(ss.hash, rand.Reader, privateKey, msg.Data, label)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(plainText, &data); err != nil {
		return err
	}

	symmetricKey, ok := data["symmetric-key"]
	if !ok {
		return errors.New("Symmetric key doesn't exist")
	}

	ss.st.Set(clientID.(string), symmetricKey.(string), 10 * time.Minute)

	return nil
}
