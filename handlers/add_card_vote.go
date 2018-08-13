package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gobuffalo/uuid"
	"github.com/nicolaszein/go-retro/models"
)

func (e Env) AddCardVote(w http.ResponseWriter, r *http.Request) {
	response := DefaultResponse{}
	card := models.Card{}

	defer r.Body.Close()

	cardID := chi.URLParam(r, "cardID")
	cardUUID, err := uuid.FromString(cardID)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.AddError("card_id", "invalid format")
		response.Write(w)
		return
	}

	if err := e.DB.FetchCardByID(cardUUID, &card); err != nil {
		response.Code = http.StatusNotFound
		response.AddError("card_id", "card does not exist")
		response.Write(w)
		return
	}

	retrospectiveID := chi.URLParam(r, "retrospectiveID")
	retrospectiveUUID, err := uuid.FromString(retrospectiveID)
	fmt.Println("handler: ", retrospectiveUUID)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.AddError("retrospective_id", "invalid format")
		response.Write(w)
		return
	}
	if retrospectiveUUID != card.RetrospectiveID {
		response.Code = http.StatusBadRequest
		response.AddError("card_id", "card does not belong to given retrospective")
		response.Write(w)
		return
	}

	card.AddVote()

	if err := e.DB.Save(&card); err != nil {
		response.Code = http.StatusInternalServerError
		response.AddError("db", "error trying to save card")
		response.Write(w)
		return
	}

	response.Code = http.StatusOK
	response.Data = card
	response.Write(w)
}
