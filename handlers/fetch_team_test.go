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

func TestFetchTeamHandler(t *testing.T) {
	testDB.CleanDatabase()
	type response struct {
		Data   models.Team       `json:"data"`
		Errors map[string]string `json:"errors"`
	}
	rctx := chi.NewRouteContext()

	t.Run("with persisted team", func(t *testing.T) {
		testDB.CleanDatabase()
		env = Env{
			DB: testDB,
		}
		team := models.Team{Name: "Team Bacon"}
		if err := testDB.Create(&team); err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/teams/{teamID}", nil)
		rctx.URLParams.Add("teamID", team.ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		env.FetchTeam(rr, r)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		if res.Data.ID != team.ID {
			t.Fatalf("team_id should be %v, but got %v", team.ID.String(), res.Data.ID)
		}
	})

	t.Run("with invalid team uuid", func(t *testing.T) {
		testDB.CleanDatabase()
		env = Env{
			DB: testDB,
		}
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/teams/{teamID}", nil)
		rctx.URLParams.Add("teamID", "invalid")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		env.FetchTeam(rr, r)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		if res.Errors["team_id"] != "invalid format" {
			t.Fatalf("team_id error should be invalid format, but got %v", res.Errors["team_id"])
		}
	})

	t.Run("with nonexistent team", func(t *testing.T) {
		testDB.CleanDatabase()
		env = Env{
			DB: testDB,
		}
		rr := httptest.NewRecorder()
		uuid, _ := uuid.NewV4()
		r := httptest.NewRequest("POST", "/api/v1/teams/{teamID}/cards", nil)
		rctx.URLParams.Add("teamID", uuid.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		env.FetchTeam(rr, r)

		if status := rr.Code; status != http.StatusNotFound {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		if res.Errors["team_id"] != "team does not exist" {
			t.Fatalf("team_id error should be team does not exist, but got %v", res.Errors["team_id"])
		}
	})
}
