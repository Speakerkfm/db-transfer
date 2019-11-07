package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"db-transfer/producer/pkg/models"
	"encoding/hex"
	"encoding/json"
	"github.com/mkideal/pkg/netutil/protocol"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net"
)

const (
	ImportData    = "IMPORT_DATA"
	createSession = "CREATE_SESSION"
)

type TcpClient struct {
	address string
	conn *Conn
}

func NewTcpClient(address string, conn *Conn) *TcpClient {
	return &TcpClient{address: address, conn: conn}
}

func (c *TcpClient) SendTcpRequest(data []byte) error {
	var err error

	c.conn.Conn, err = net.Dial(protocol.TCP, c.address)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(append(data, '\n'))
	c.conn.Close()

	return err
}

func (c *TcpClient) SendRSARequest(clientID string, publicKey *rsa.PublicKey, msgType string, data []byte) error {
	headers := make(map[string]interface{}, 2)
	headers["X-MESSAGE-TYPE"] = msgType
	headers["X-CLIENT-ID"] = clientID

	label := []byte("")
	hash := sha256.New()

	encryptedData, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, data, label)
	if err != nil {
		return err
	}

	msg := &models.Message{
		Headers: headers,
		Data:    encryptedData,
	}

	reqBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c.conn.Conn, err = net.Dial(protocol.TCP, c.address)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(append(reqBody, '\n'))
	c.conn.Close()

	return err
}

func (c *TcpClient) RequestPublicKey() (*rsa.PublicKey, error) {
	emptyMsg, err := json.Marshal(models.Message{})
	if err != nil {
		return nil, err
	}

	c.conn.Conn, err = net.Dial(protocol.TCP, c.address)
	if err != nil {
		return nil, err
	}

	c.conn.Write(append(emptyMsg, '\n'))

	respBody, err := ioutil.ReadAll(c.conn)
	if err != nil {
		panic(err)
	}

	c.conn.Close()

	var msg *models.Message
	if err := json.Unmarshal(respBody, &msg); err != nil {
		return nil, err
	}

	var data map[string]rsa.PublicKey
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return nil, err
	}


	publicKey, ok := data["public-key"]
	if !ok {
		return nil, errors.New("Symmetric key doesn't exist")
	}

	return &publicKey, nil
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
