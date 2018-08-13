package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nicolaszein/go-retro/models"
)

func TestListTeamsHandler(t *testing.T) {
	testDB.CleanDatabase()
	env = Env{
		DB: testDB,
	}
	handler := http.HandlerFunc(env.ListTeams)
	type response struct {
		Data   []models.Team     `json:"data"`
		Errors map[string]string `json:"errors"`
	}

	t.Run("with persisteds teams", func(t *testing.T) {
		testDB.CleanDatabase()
		team := models.Team{Name: "Team Bacon"}
		if err := testDB.Create(&team); err != nil {
			t.Fatalf("Failed creating team %v", err)
		}
		team2 := models.Team{Name: "Team Bacon2"}
		if err := testDB.Create(&team2); err != nil {
			t.Fatalf("Failed creating team2 %v", err)
		}
		res := response{}
		req, err := http.NewRequest("GET", "/api/v1/teams", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		if len(res.Data) != 2 {
			t.Fatalf("should return 2 teams, but got %v", len(res.Data))
		}
	})

	t.Run("with no teams", func(t *testing.T) {
		testDB.CleanDatabase()
		res := response{}
		req, err := http.NewRequest("GET", "/api/v1/teams", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		if len(res.Data) != 0 {
			t.Fatalf("should return 0 teams, but got %v", len(res.Data))
		}
	})

	t.Run("with db error", func(t *testing.T) {
		dbMock.Error = errors.New("Error trying to fetch teams")
		env.DB = dbMock
		handler := http.HandlerFunc(env.ListTeams)
		res := response{}
		req, err := http.NewRequest("GET", "/api/v1/teams", nil)
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
		if msg := "Error trying to fetch teams"; res.Errors["db"] != msg {
			t.Fatalf("handler should return db error: %s, but got %v", msg, res.Errors["db"])
		}
	})
}
