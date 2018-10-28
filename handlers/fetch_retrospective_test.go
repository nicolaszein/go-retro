package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/gobuffalo/uuid"
	"github.com/nicolaszein/go-retro/models"
)

func TestFetchRetrospectiveHandler(t *testing.T) {
	testDB.CleanDatabase()
	type response struct {
		Data   models.Retrospective `json:"data"`
		Errors map[string]string    `json:"errors"`
	}
	rctx := chi.NewRouteContext()

	t.Run("with persisted retrospective", func(t *testing.T) {
		testDB.CleanDatabase()
		env = Env{
			DB: testDB,
		}
		team := models.Team{Name: "Team Bacon"}
		retrospective := models.Retrospective{Name: "Retrospective Bacon", Team: team}
		if err := testDB.Create(&retrospective); err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/retrospectives/{retrospectiveID}", nil)
		rctx.URLParams.Add("retrospectiveID", retrospective.ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		env.FetchRetrospective(rr, r)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		if res.Data.ID != retrospective.ID {
			t.Fatalf("retrospective_id should be %v, but got %v", retrospective.ID.String(), res.Data.ID)
		}
	})

	t.Run("with invalid retrospective uuid", func(t *testing.T) {
		testDB.CleanDatabase()
		env = Env{
			DB: testDB,
		}
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/retrospectives/{retrospectiveID}", nil)
		rctx.URLParams.Add("retrospectiveID", "invalid")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		env.FetchRetrospective(rr, r)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		if res.Errors["retrospective_id"] != "invalid format" {
			t.Fatalf("retrospective_id error should be invalid format, but got %v", res.Errors["retrospective_id"])
		}
	})

	t.Run("with nonexistent retrospective", func(t *testing.T) {
		testDB.CleanDatabase()
		env = Env{
			DB: testDB,
		}
		rr := httptest.NewRecorder()
		uuid, _ := uuid.NewV4()
		r := httptest.NewRequest("POST", "/api/v1/retrospectives/{retrospectiveID}/cards", nil)
		rctx.URLParams.Add("retrospectiveID", uuid.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		env.FetchRetrospective(rr, r)

		if status := rr.Code; status != http.StatusNotFound {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		if res.Errors["retrospective_id"] != "retrospective does not exist" {
			t.Fatalf("retrospective_id error should be retrospective does not exist, but got %v", res.Errors["restrospective_id"])
		}
	})
}
