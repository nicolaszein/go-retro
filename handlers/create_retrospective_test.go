package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gobuffalo/uuid"
	"github.com/nicolaszein/go-retro/models"
)

func TestCreateRetrospectiveHandler(t *testing.T) {
	testDB.CleanDatabase()
	env = Env{
		DB: testDB,
	}
	handler := http.HandlerFunc(env.CreateRetrospective)
	type response struct {
		Data   models.Retrospective `json:"data"`
		Errors map[string]string    `json:"errors"`
	}
	t.Run("with invalid params", func(t *testing.T) {
		testDB.CleanDatabase()
		params := strings.NewReader(`{}`)
		req, err := http.NewRequest("POST", "/api/v1/retrospectives", params)
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

	t.Run("with invalid body", func(t *testing.T) {
		testDB.CleanDatabase()
		res := response{}
		params := strings.NewReader(``)
		req, err := http.NewRequest("POST", "/api/v1/retrospectives", params)
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
			t.Fatalf("body error should be invalid format, but got %v", res.Errors["body"])
		}
	})

	t.Run("with nonexistent team", func(t *testing.T) {
		testDB.CleanDatabase()
		res := response{}
		uuid, err := uuid.NewV4()
		payload := fmt.Sprintf(`{"name": "Retrospective", "team_id": "%v"}`, uuid)
		params := strings.NewReader(payload)
		req, err := http.NewRequest("POST", "/api/v1/retrospectives", params)
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

		if res.Errors["team_id"] != "invalid team_id" {
			t.Fatalf("team_id error should be invalid team_id, but got %v", res.Errors["team_id"])
		}
	})

	t.Run("with valid params", func(t *testing.T) {
		testDB.CleanDatabase()
		team := models.Team{Name: "Team A"}
		if err := testDB.Create(&team); err != nil {
			t.Fatalf("Error trying to create team! Error: %v", err)
		}
		res := response{}
		payload := fmt.Sprintf(`{"name": "Retrospective", "team_id": "%v"}`, team.ID)
		params := strings.NewReader(payload)
		req, err := http.NewRequest("POST", "/api/v1/retrospectives", params)
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

		if res.Data.Name != "Retrospective" {
			t.Fatalf("restrospective name should be Retropective, but got %v", res.Data.Name)
		}

		if res.Data.TeamID != team.ID {
			t.Fatalf("restrospective team_id should be %v, but got %v", team.ID, res.Data.TeamID)
		}
	})

	t.Run("with db error", func(t *testing.T) {
		dbMock.Error = errors.New("error trying to create retrospective")
		env = Env{
			DB: dbMock,
		}
		handler := http.HandlerFunc(env.CreateRetrospective)
		res := response{}
		uuid, err := uuid.NewV4()
		payload := fmt.Sprintf(`{"name": "Retrospective", "team_id": "%v"}`, uuid)
		params := strings.NewReader(payload)
		req, err := http.NewRequest("POST", "/api/v1/retrospectives", params)
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

		if msg := "error trying to create retrospective"; res.Errors["db"] != msg {
			t.Fatalf("handler should return db error: %s, but got %v", msg, res.Errors["db"])
		}
	})
}
