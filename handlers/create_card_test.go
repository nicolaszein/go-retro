package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/gobuffalo/uuid"
	"github.com/nicolaszein/go-retro/models"
)

func TestCreateCardHandler(t *testing.T) {
	type response struct {
		Data   models.Card       `json:"data"`
		Errors map[string]string `json:"errors"`
	}
	rctx := chi.NewRouteContext()

	t.Run("with nonexistent retrospective", func(t *testing.T) {
		testDB.CleanDatabase()
		env = Env{
			DB: testDB,
		}
		params := strings.NewReader(`{}`)
		rr := httptest.NewRecorder()
		uuid, _ := uuid.NewV4()
		r := httptest.NewRequest("POST", "/api/v1/retrospectives/{retrospectiveID}/cards", params)
		rctx.URLParams.Add("retrospectiveID", uuid.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		env.CreateCard(rr, r)

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

	t.Run("with invalid retrospective uuid", func(t *testing.T) {
		testDB.CleanDatabase()
		env = Env{
			DB: testDB,
		}
		params := strings.NewReader(`{}`)
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/retrospectives/{retrospectiveID}/cards", params)
		rctx.URLParams.Add("retrospectiveID", "invalid")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		env.CreateCard(rr, r)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		if res.Errors["retrospective_id"] != "invalid format" {
			t.Fatalf("retrospective_id error should be invalid format, but got %v", res.Errors["restrospective_id"])
		}
	})

	t.Run("with invalid body", func(t *testing.T) {
		testDB.CleanDatabase()
		env = Env{
			DB: testDB,
		}
		team := models.Team{Name: "Team Bacon"}
		retrospective := models.Retrospective{Name: "Retro bacon", Team: team}
		if err := testDB.Create(&retrospective); err != nil {
			t.Fatal(err)
		}
		params := strings.NewReader(``)
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/retrospectives/{retrospectiveID}/cards", params)
		rctx.URLParams.Add("retrospectiveID", retrospective.ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		env.CreateCard(rr, r)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if res.Errors["body"] != "invalid format" {
			t.Fatalf("body error should be invalid format, but got %v", res.Errors["body"])
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

		env.CreateCard(rr, r)

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

		env.CreateCard(rr, r)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if msg := "invalid does not validate as matches(^positive$|^negative$)"; res.Errors["type"] != msg {
			t.Fatalf("type error should be %v, but got %v", msg, res.Errors["type"])
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

		env.CreateCard(rr, r)

		if status := rr.Code; status != http.StatusCreated {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}
	})

	t.Run("with db error", func(t *testing.T) {
		dbMock.CleanDatabase()
		dbMock.Error = errors.New("error trying to create card")
		env = Env{
			DB: dbMock,
		}
		uuid, _ := uuid.NewV4()
		params := strings.NewReader(`{"content": "card content", "type": "positive"}`)
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/retrospectives/{retrospectiveID}/cards", params)
		rctx.URLParams.Add("retrospectiveID", uuid.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		env.CreateCard(rr, r)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if msg := "error trying to create retrospective"; res.Errors["db"] != msg {
			t.Fatalf("handler should return db error: %s, but got %v", msg, res.Errors["db"])
		}
	})
}
