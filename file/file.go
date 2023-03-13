package file

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

var fileMutex sync.Mutex

type savedFile map[string]string

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

func SaveSecrets(f string, d map[string]string) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	bytes, err := json.Marshal(d)
	if err != nil {
		fmt.Println("error encoding data", err)
	}
	err = os.WriteFile(f, bytes, 0666)
	if err != nil {
		fmt.Print("Error writing to file", err)
	}
}
