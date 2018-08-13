package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nicolaszein/go-retro/models"
)

func TestCreateTeamHandler(t *testing.T) {
	testDB.CleanDatabase()
	env = Env{
		DB: testDB,
	}
	handler := http.HandlerFunc(env.CreateTeam)
	type response struct {
		Data   models.Team       `json:"data"`
		Errors map[string]string `json:"errors"`
	}

	t.Run("with valid payload", func(t *testing.T) {
		testDB.CleanDatabase()
		res := response{}
		params := strings.NewReader(`{"name": "Team Bacon"}`)
		req, err := http.NewRequest("POST", "/api/v1/teams", params)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		if res.Data.Name != "Team Bacon" {
			t.Fatalf("team name should be Team Bacon but got %v", res.Data.Name)
		}
	})

	t.Run("with invalid payload", func(t *testing.T) {
		params := strings.NewReader(`{}`)
		req, err := http.NewRequest("POST", "/api/v1/teams", params)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("with invalid body format", func(t *testing.T) {
		res := response{}
		params := strings.NewReader(``)
		req, err := http.NewRequest("POST", "/api/v1/teams", params)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}

		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if res.Errors["body"] != "invalid format" {
			t.Fatalf("handler should return invalid body format, but got %v", res.Errors["body"])
		}
	})

	t.Run("with db error", func(t *testing.T) {
		dbMock.Error = errors.New("Fuuuu")
		env.DB = dbMock
		handler := http.HandlerFunc(env.CreateTeam)

		res := response{}
		params := strings.NewReader(`{"name": "Team Bacon"}`)
		req, err := http.NewRequest("POST", "/api/v1/teams", params)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}

		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if msg := "error trying to create team"; res.Errors["db"] != msg {
			t.Fatalf("handler should return db error: %s, but got %v", msg, res.Errors["db"])
		}
	})
}
