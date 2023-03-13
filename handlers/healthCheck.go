package handlers

import (
	"fmt"
	"io"
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got /healthcheck request")
	io.WriteString(w, "Ok")
}
