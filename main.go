package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nicolaszein/go-retro/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println("Starting server on port " + port)

	http.HandleFunc("/", handlers.HealthCheck)
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		fmt.Println("Error serving:", err)
	}
}
