package main

import (
	"fmt"
	"net/http"

	"github.com/nicolaszein/go-retro/handlers"
	"github.com/nicolaszein/go-retro/settings"
)

func main() {
	port := settings.PORT
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
