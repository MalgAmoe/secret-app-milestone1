package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"secret-app/file"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type secret struct {
	PlainText string `json:"plain_text"`
}

type secretId struct {
	Id string `json:"id"`
}

type secretResponse struct {
	Data string `json:"data"`
}

func getSecretHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got / GET request")

	decoder := json.NewDecoder(r.Body)
	var id secretId
	err := decoder.Decode(&id)
	if err != nil {
		fmt.Println("JSON decoding error", err)
		w.WriteHeader(http.StatusNotFound)
		response := secretResponse{
			"",
		}
		bytes, _ := json.Marshal(response)
		io.WriteString(w, string(bytes))
		return
	}

	secretString, err := file.DaFile.RemoveSecret(id.Id)
	if err != nil || secretString == "" {
		w.WriteHeader(http.StatusNotFound)
		response := secretResponse{
			"",
		}
		bytes, _ := json.Marshal(response)
		io.WriteString(w, string(bytes))
		return
	}

	response := secretResponse{
		secretString,
	}
	bytes, _ := json.Marshal(response)

	io.WriteString(w, string(bytes))
}

func postSecretHandler(w http.ResponseWriter, r *http.Request /* , d map[string]string, c func() */) {
	fmt.Println("got / POST request")

	decoder := json.NewDecoder(r.Body)
	var secret secret
	err := decoder.Decode(&secret)
	if err != nil {
		fmt.Println("JSON decoding error", err)
		io.WriteString(w, "error decoding the body")
		return
	}

	h := file.GetMD5Hash(secret.PlainText)
	file.DaFile.AddSecret(secret.PlainText, h)
	response := secretId{
		h,
	}
	bytes, _ := json.Marshal(response)

	io.WriteString(w, string(bytes))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got /healthcheck request")
	io.WriteString(w, "Ok")
}

func main() {
	file.Init()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getSecretHandler(w, r)
		} else if r.Method == http.MethodPost {
			postSecretHandler(w, r)
		}
	})
	mux.HandleFunc("/healthcheck", healthCheckHandler)

	err := http.ListenAndServe(":8080", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(0)
	}
}
