package main

import (
	"fmt"
	"net/http"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ok")
}

func main() {
	fmt.Println("Starting server on port :8000")

	http.HandleFunc("/", healthCheck)
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		fmt.Println("Error serving:", err)
	}
}
