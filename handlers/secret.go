package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"secret-app/file"
)

type secret struct {
	PlainText string `json:"plain_text"`
}

type secretId struct {
	Id string `json:"id"`
}

type secretResponse struct {
	Data string `json:"data"`
}

func secretHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getSecretHandler(w, r)
	case http.MethodPost:
		postSecretHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
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

func postSecretHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got / POST request")

	decoder := json.NewDecoder(r.Body)
	var secret secret
	err := decoder.Decode(&secret)
	if err != nil {
		fmt.Println("JSON decoding error", err)
		response := secretId{
			"",
		}
		bytes, _ := json.Marshal(response)
		io.WriteString(w, string(bytes))
		return
	}

	h := file.GetMD5Hash(secret.PlainText)
	err = file.DaFile.AddSecret(secret.PlainText, h)
	if err != nil {
		fmt.Println("Adding secret error", err)
		response := secretId{
			"",
		}
		bytes, _ := json.Marshal(response)
		io.WriteString(w, string(bytes))
		return
	}
	response := secretId{
		h,
	}
	bytes, _ := json.Marshal(response)

	io.WriteString(w, string(bytes))
}
