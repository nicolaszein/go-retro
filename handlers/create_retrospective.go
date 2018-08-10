package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"

	"github.com/nicolaszein/go-retro/models"
)

func (e Env) CreateRetrospective(w http.ResponseWriter, r *http.Request) {
	response := DefaultResponse{}
	retrospective := models.Retrospective{}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&retrospective); err != nil {
		response.Code = http.StatusBadRequest
		response.AddError("body", "invalid format")
		response.Write(w)
		return
	}

	_, err := govalidator.ValidateStruct(retrospective)
	if err != nil {
		errors := govalidator.ErrorsByField(err)
		response.Code = http.StatusBadRequest
		response.Errors = errors
		response.Write(w)
		return
	}

	if err := e.DB.FetchTeamByID(retrospective.TeamID, &retrospective.Team); err != nil {
		response.Code = http.StatusBadRequest
		response.AddError("team_id", "invalid team_id")
		response.Write(w)
		return
	}

	if err := e.DB.Create(&retrospective); err != nil {
		response.Code = http.StatusInternalServerError
		response.AddError("db", "error trying to create retrospective")
		response.Write(w)
		return
	}

	response.Code = http.StatusCreated
	response.Data = retrospective
	response.Write(w)
}
