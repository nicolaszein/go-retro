package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gobuffalo/uuid"
	"github.com/nicolaszein/go-retro/models"
)

func (e Env) FetchTeam(w http.ResponseWriter, r *http.Request) {
	response := DefaultResponse{}
	team := models.Team{}

	defer r.Body.Close()

	teamID := chi.URLParam(r, "teamID")
	teamUUID, err := uuid.FromString(teamID)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.AddError("team_id", "invalid format")
		response.Write(w)
		return
	}

	if err := e.DB.FetchTeamByID(teamUUID, &team); err != nil {
		response.Code = http.StatusNotFound
		response.AddError("team_id", "team does not exist")
		response.Write(w)
		return
	}

	response.Code = http.StatusOK
	response.Data = team
	response.Write(w)
}
