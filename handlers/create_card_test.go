package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/nicolaszein/go-retro/models"
)

func TestCreateCardHandler(t *testing.T) {
	type response struct {
		Data   models.Card       `json:"data"`
		Errors map[string]string `json:"errors"`
	}
	handler := env.CreateCard
	rctx := chi.NewRouteContext()

	t.Run("with nonexistent retrospective", func(t *testing.T) {
		testDB.CleanDatabase()
		params := strings.NewReader(`{}`)
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/retrospectives/{retrospectiveID}/cards", params)
		rctx.URLParams.Add("retrospectiveID", "value")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		handler(rr, r)

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

	t.Run("with invalid params", func(t *testing.T) {
		testDB.CleanDatabase()
		env = Env{
			DB: testDB,
		}
		team := models.Team{Name: "Team Bacon"}
		retrospective := models.Retrospective{Name: "Retro bacon", Team: team}
		if err := testDB.Create(&retrospective); err != nil {
			t.Fatal(err)
		}
		params := strings.NewReader(`{}`)
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/retrospectives/{retrospectiveID}/cards", params)
		rctx.URLParams.Add("retrospectiveID", retrospective.ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		handler(rr, r)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("with invalid type", func(t *testing.T) {
		testDB.CleanDatabase()
		env = Env{
			DB: testDB,
		}
		team := models.Team{Name: "Team Bacon"}
		retrospective := models.Retrospective{Name: "Retro bacon", Team: team}
		if err := testDB.Create(&retrospective); err != nil {
			t.Fatal(err)
		}
		params := strings.NewReader(`{"type": "invalid"}`)
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/retrospectives/{retrospectiveID}/cards", params)
		rctx.URLParams.Add("retrospectiveID", retrospective.ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		handler(rr, r)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if msg := "xablau"; res.Errors["type"] != msg {
			t.Fatalf("type error should be %v, but got %v", msg, res.Errors["restrospective_id"])
		}
	})

	t.Run("with valid params", func(t *testing.T) {
		testDB.CleanDatabase()
		env = Env{
			DB: testDB,
		}
		team := models.Team{Name: "Team Bacon"}
		retrospective := models.Retrospective{Name: "Retro bacon", Team: team}
		if err := testDB.Create(&retrospective); err != nil {
			t.Fatal(err)
		}
		params := strings.NewReader(`{"content": "card content", "type": "positive"}`)
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/retrospectives/{retrospectiveID}/cards", params)
		rctx.URLParams.Add("retrospectiveID", retrospective.ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		handler(rr, r)

		if status := rr.Code; status != http.StatusCreated {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}
	})
}
