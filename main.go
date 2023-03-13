package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"secret-app/file"
	"secret-app/handlers"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	file.Init()

	mux := http.NewServeMux()
	handlers.SetupHandlers(mux)

	err := http.ListenAndServe(":8080", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(0)
	}
}
