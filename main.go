package main

import (
	"fmt"
	"net/http"

	"github.com/nicolaszein/go-retro/handlers"
)

func main() {
	fmt.Println("Starting server on port :8000")

	http.HandleFunc("/", handlers.HealthCheck)
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		fmt.Println("Error serving:", err)
	}
}
