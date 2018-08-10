package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi"
	"github.com/gobuffalo/uuid"
	"github.com/nicolaszein/go-retro/models"
)

func (e Env) CreateCard(w http.ResponseWriter, r *http.Request) {
	response := DefaultResponse{}
	card := models.Card{}

	defer r.Body.Close()

	retrospectiveID := chi.URLParam(r, "retrospectiveID")
	retrospectiveUUID, err := uuid.FromString(retrospectiveID)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.AddError("retrospective_id", "invalid format")
		response.Write(w)
		return
	}

	if err := e.DB.FetchRetrospectiveByID(retrospectiveUUID, &card.Retrospective); err != nil {
		response.Code = http.StatusNotFound
		response.AddError("retrospective_id", "retrospective does not exist")
		response.Write(w)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		response.Code = http.StatusBadRequest
		response.AddError("body", "invalid format")
		response.Write(w)
		return
	}

	_, err = govalidator.ValidateStruct(card)
	if err != nil {
		errors := govalidator.ErrorsByField(err)
		response.Code = http.StatusBadRequest
		response.Errors = errors
		response.Write(w)
		return
	}

	if err := e.DB.Create(&card); err != nil {
		response.Code = http.StatusInternalServerError
		response.AddError("db", "error trying to create retrospective")
		response.Write(w)
		return
	}

	response.Code = http.StatusCreated
	response.Data = card
	response.Write(w)
}
