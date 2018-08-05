package handlers

import (
	"net/http"

	"github.com/nicolaszein/go-retro/models"
)

func (e Env) ListTeams(w http.ResponseWriter, r *http.Request) {
	response := DefaultResponse{}
	teams := []models.Team{}

	if err := e.DB.FetchTeams(&teams); err != nil {
		response.Code = http.StatusInternalServerError
		response.AddError("db", "Error trying to fetch teams")
		response.Write(w)
		return
	}

	response.Data = teams
	response.Code = http.StatusOK
	response.Write(w)
}
