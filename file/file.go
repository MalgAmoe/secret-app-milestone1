package file

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

var fileMutex sync.Mutex

type savedFile map[string]string

type fileData struct {
	path string
	file savedFile
}

var DaFile fileData

func Init() {
	f := CreateFile()
	DaFile = fileData{f, GetData(f)}
}

func CreateFile() string {
	fPath := os.Getenv("DATA_FILE_PATH")
	if fPath == "" {
		log.Fatal("DATA_FILE_PATH is not set")
	}
	_, err := os.Stat(fPath)
	if err != nil {
		_, err = os.Create(fPath)
		if err != nil {
			log.Fatal("Could not create the data file")
		}
	}
	return fPath
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
	fileMutex.Lock()
	defer fileMutex.Unlock()

	f.file[h] = s
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
	fileMutex.Lock()
	defer fileMutex.Unlock()

	s := f.file[h]
	if s == "" {
		return s, nil
	}
	delete(f.file, h)

	bytes, err := json.Marshal(f.file)
	if err != nil {
		fmt.Println("error encoding data", err)
	}

	// TODO: possible internal error here
	os.WriteFile(f.path, bytes, 0666)

	return s, nil
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
