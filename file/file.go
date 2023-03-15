package file

import (
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"secret-app/crypto"
)

type savedFile map[string]string

type fileData struct {
	path  string
	file  savedFile
	Mu    sync.Mutex
	GCM   cipher.AEAD
	Nonce []byte
}

var DaFile fileData

func Init() {
	f, gcm, nonce := CreateFile()
	DaFile = fileData{f, GetData(f), sync.Mutex{}, gcm, nonce}
}

func CreateFile() (string, cipher.AEAD, []byte) {
	fPath := os.Getenv("DATA_FILE_PATH")
	password := os.Getenv("PASSWORD")
	salt := os.Getenv("SALT")
	if fPath == "" {
		log.Fatal("DATA_FILE_PATH is not set")
	}
	if password == "" {
		log.Fatal("PASSWORD is not set")
	}
	if salt == "" {
		log.Fatal("SALT is not set")
	}
	gcm, nonce, err := crypto.InitCrypto(password, salt)
	if err != nil {
		log.Fatal("could not initialize crypto primitives")
	}
	_, err = os.Stat(fPath)
	if err != nil {
		_, err = os.Create(fPath)
		if err != nil {
			log.Fatal("Could not create the data file")
		}
	}
	return fPath, gcm, nonce
}

func GetData(fPath string) savedFile {
	fRead, err := os.ReadFile(fPath)
	if err != nil {
		log.Fatal(err)
	}
	var dataFile savedFile

	err = json.Unmarshal(fRead, &dataFile)
	if err != nil {
		return make(map[string]string)
	} else {
		return dataFile
	}
}

func (f *fileData) AddSecret(s string, h string) error {
	f.Mu.Lock()
	defer f.Mu.Unlock()

	encData := crypto.Encrypt(s, f.GCM, f.Nonce)
	f.file[h] = string(encData)
	bytes, err := json.Marshal(f.file)
	if err != nil {
		fmt.Println("error encoding data", err)
		return err
	}

	err = os.WriteFile(f.path, bytes, 0666)
	if err != nil {
		fmt.Print("Error writing to file", err)
		return err
	}

	return nil
}

func (f *fileData) RemoveSecret(h string) (string, error) {
	f.Mu.Lock()
	defer f.Mu.Unlock()

	s := f.file[h]
	if s == "" {
		return s, nil
	}
	delete(f.file, h)

	bytes, err := json.Marshal(f.file)
	if err != nil {
		fmt.Println("error encrypt data", err)
	}

	// TODO: possible internal error here
	os.WriteFile(f.path, bytes, 0666)

	decData, err := crypto.Decrypt([]byte(s), f.GCM)
	if err != nil {
		fmt.Println("Could not decrypt data", err)
	}

	return string(decData), nil
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
