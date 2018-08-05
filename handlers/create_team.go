package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nicolaszein/go-retro/models"

	"github.com/asaskevich/govalidator"
)

func (e Env) CreateTeam(w http.ResponseWriter, r *http.Request) {
	response := DefaultResponse{}
	team := models.Team{}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		response.Code = http.StatusBadRequest
		response.AddError("body", "invalid format")
		response.Write(w)
		return
	}

	_, err := govalidator.ValidateStruct(team)
	if err != nil {
		errors := govalidator.ErrorsByField(err)
		response.Code = http.StatusBadRequest
		response.Errors = errors
		response.Write(w)
		return
	}

	if err := e.DB.Create(&team); err != nil {
		response.Code = http.StatusInternalServerError
		response.AddError("db", "error trying to create team")
		response.Write(w)
		return
	}

	response.Code = http.StatusCreated
	response.Data = team
	response.Write(w)
}
