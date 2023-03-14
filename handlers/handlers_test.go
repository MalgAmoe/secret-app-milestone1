package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"secret-app/file"
	"strings"
	"testing"
)

var mux *http.ServeMux
var writer *httptest.ResponseRecorder

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	os.Exit(code)
}

func setUp() {
	mux = http.NewServeMux()
	mux.HandleFunc("/healthcheck", healthCheckHandler)
	mux.HandleFunc("/", secretHandler)

	os.Setenv("DATA_FILE_PATH", "testMocks.json")
	defer os.Unsetenv("DATA_FILE_PATH")
	file.Init()
}

func TestHealthCheckReguest(t *testing.T) {
	request, _ := http.NewRequest("GET", "/healthcheck", nil)
	writer = httptest.NewRecorder()
	mux.ServeHTTP(writer, request)

	if writer.Code != http.StatusOK {
		t.Errorf("Response code %v", writer.Code)
	}

	if writer.Body.String() != "Ok" {
		t.Error("healthcheck is not Ok")
	}
}

func TestSecretRequestMethods(t *testing.T) {
	json := strings.NewReader(`{}`)
	request, _ := http.NewRequest("GET", "/", json)
	writer = httptest.NewRecorder()
	mux.ServeHTTP(writer, request)
	if writer.Code == http.StatusMethodNotAllowed {
		t.Error("GET method is not accepted")
	}

	json = strings.NewReader(`{}`)
	request, _ = http.NewRequest("POST", "/", json)
	writer = httptest.NewRecorder()
	mux.ServeHTTP(writer, request)
	if writer.Code == http.StatusMethodNotAllowed {
		t.Error("POST method is not accepted")
	}

	json = strings.NewReader(`{}`)
	request, _ = http.NewRequest("PUT", "/", json)
	writer = httptest.NewRecorder()
	mux.ServeHTTP(writer, request)
	if writer.Code != http.StatusMethodNotAllowed {
		t.Error("PUT method is allowed and should not")
	}
}

func TestSecretCreationAndRetrieval(t *testing.T) {
	jsonBody := strings.NewReader(`{ "plain_text": "secret" }`)
	request, _ := http.NewRequest("POST", "/", jsonBody)
	writer = httptest.NewRecorder()
	mux.ServeHTTP(writer, request)
	if writer.Code != http.StatusOK {
		t.Error("POST method for adding secret failed")
	}

	jsonBody = strings.NewReader(`{ "id": "5ebe2294ecd0e0f08eab7690d2a6ee69" }`)
	request, _ = http.NewRequest("GET", "/", jsonBody)
	writer = httptest.NewRecorder()
	mux.ServeHTTP(writer, request)
	if writer.Code != http.StatusOK {
		t.Error("GET method to retrieve secret failed")
	}
	var secretResponse SecretResponse
	json.Unmarshal(writer.Body.Bytes(), &secretResponse)
	if secretResponse.Data != "secret" {
		t.Error("Secret not retrieved correctly")
	}
}

func TestWrongRetrieval(t *testing.T) {
	jsonBody := strings.NewReader(`{ "id": "random" }`)
	request, _ := http.NewRequest("GET", "/", jsonBody)
	writer = httptest.NewRecorder()
	mux.ServeHTTP(writer, request)
	if writer.Code != http.StatusNotFound {
		t.Error("GET method returns wrong method")
	}
	var secretResponse SecretResponse
	json.Unmarshal(writer.Body.Bytes(), &secretResponse)
	if secretResponse.Data != "" {
		t.Error("Secret is not empty")
	}
}
