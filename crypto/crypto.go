package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"

	"golang.org/x/crypto/scrypt"
)

func InitCrypto(password, salt string) (cipher.AEAD, []byte, error) {
	var gcm cipher.AEAD
	var nonce []byte
	key, err := scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
	if err != nil {
		return gcm, nonce, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return gcm, nonce, err
	}
	gcm, err = cipher.NewGCM(block)
	if err != nil {
		return gcm, nonce, err
	}
	nonce = make([]byte, gcm.NonceSize())
	if err != nil {
		return gcm, nonce, err
	}
	return gcm, nonce, nil
}

func Encrypt(data string, gcm cipher.AEAD, nonce []byte) []byte {
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}
	encryptedData := gcm.Seal(nonce, nonce, []byte(data), nil)
	return encryptedData
}

func Decrypt(encData []byte, gcm cipher.AEAD) ([]byte, error) {
	nonce := encData[:gcm.NonceSize()]
	encData = encData[gcm.NonceSize():]
	data, err := gcm.Open(nil, nonce, encData, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}
