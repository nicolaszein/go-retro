package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func (e Env) CreateCard(w http.ResponseWriter, r *http.Request) {
	retrospectiveID := chi.URLParam(r, "retrospectiveID")
	fmt.Println(retrospectiveID)
}
