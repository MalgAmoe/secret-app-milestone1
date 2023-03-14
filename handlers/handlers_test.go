package handlers

import (
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
