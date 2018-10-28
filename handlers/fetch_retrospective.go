package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gobuffalo/uuid"
	"github.com/nicolaszein/go-retro/models"
)

func (e Env) FetchRetrospective(w http.ResponseWriter, r *http.Request) {
	response := DefaultResponse{}
	retrospective := models.Retrospective{}

	defer r.Body.Close()

	retrospectiveID := chi.URLParam(r, "retrospectiveID")
	retrospectiveUUID, err := uuid.FromString(retrospectiveID)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.AddError("retrospective_id", "invalid format")
		response.Write(w)
		return
	}

	if err := e.DB.FetchRetrospectiveByID(retrospectiveUUID, &retrospective); err != nil {
		response.Code = http.StatusNotFound
		response.AddError("retrospective_id", "retrospective does not exist")
		response.Write(w)
		return
	}

	response.Code = http.StatusOK
	response.Data = retrospective
	response.Write(w)
}
