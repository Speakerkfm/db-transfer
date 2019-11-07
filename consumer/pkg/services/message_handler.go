package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"db-transfer/consumer/pkg/models"
	"db-transfer/consumer/pkg/store"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
)

const (
	importData    = "IMPORT_DATA"
	createSession = "CREATE_SESSION"
)

type MessageHandler struct {
	t          *Transfer
	ss         *SessionService
	st         *store.Store
	privateKey *rsa.PrivateKey
}

func NewMessageHandler(t *Transfer, ss *SessionService, st *store.Store, privateKey *rsa.PrivateKey) *MessageHandler {
	return &MessageHandler{
		t:          t,
		ss:         ss,
		st:         st,
		privateKey: privateKey,
	}
}

func (mh *MessageHandler) Handle(msg *models.Message) (*models.Message, error) {
	if msg == nil {
		return nil, nil
	}

	msgType, ok := msg.Headers["X-MESSAGE-TYPE"]
	if !ok {
		return mh.publicKeyResponse(&mh.privateKey.PublicKey), nil
	}
	switch msgType {
	case importData:
		return nil, mh.importData(msg)
	case createSession:
		return nil, mh.ss.CreateSession(msg, mh.privateKey)
	default:
		return mh.publicKeyResponse(&mh.privateKey.PublicKey), nil
	}

	return nil, nil
}

func (mh *MessageHandler) importData(msg *models.Message) error {
	clientID, ok := msg.Headers["X-CLIENT-ID"]
	if !ok {
		return errors.New("Unknown client")
	}

	symmetricKey := mh.st.Get(clientID.(string))
	decryptedData := decrypt(msg.Data, symmetricKey)

	return mh.t.Handle(decryptedData)
}

func (mh *MessageHandler) publicKeyResponse(publicKey *rsa.PublicKey) *models.Message {
	responseData := make(map[string]interface{}, 1)
	responseData["public-key"] = publicKey

	data, err := json.Marshal(responseData)
	if err != nil {
		panic(err)
	}

	return &models.Message{
		Data: data,
	}
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}
