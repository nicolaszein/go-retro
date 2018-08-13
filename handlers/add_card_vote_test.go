package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/gobuffalo/uuid"
	"github.com/nicolaszein/go-retro/models"
)

func TestAddCardVoteHandler(t *testing.T) {
	testDB.CleanDatabase()
	type response struct {
		Data   models.Card       `json:"data"`
		Errors map[string]string `json:"errors"`
	}
	rctx := chi.NewRouteContext()

	t.Run("with nonexistent card", func(t *testing.T) {
		testDB.CleanDatabase()
		env = Env{
			DB: testDB,
		}
		rr := httptest.NewRecorder()
		uuid, _ := uuid.NewV4()
		r := httptest.NewRequest("POST", "/api/v1/retrospectives/{retrospectiveID}/cards/{cardID}/votes", nil)
		rctx.URLParams.Add("retrospectiveID", uuid.String())
		rctx.URLParams.Add("cardID", uuid.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		env.AddCardVote(rr, r)

		if status := rr.Code; status != http.StatusNotFound {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		if res.Errors["card_id"] != "card does not exist" {
			t.Fatalf("card_id error should be card does not exist, but got %v", res.Errors["card_id"])
		}
	})

	t.Run("with invalid card uuid", func(t *testing.T) {
		testDB.CleanDatabase()
		env = Env{
			DB: testDB,
		}
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/retrospectives/{retrospectiveID}/cards", nil)
		rctx.URLParams.Add("cardID", "invalid")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		env.AddCardVote(rr, r)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		if res.Errors["card_id"] != "invalid format" {
			t.Fatalf("card_id error should be invalid format, but got %v", res.Errors["card_id"])
		}
	})

	t.Run("with invalid retrospective uuid", func(t *testing.T) {
		testDB.CleanDatabase()
		env = Env{
			DB: testDB,
		}
		team := models.Team{Name: "Team Bacon"}
		retrospective := models.Retrospective{Name: "Retrospective Bacon", Team: team}
		card := models.Card{Content: "Card Bacon", Type: "negative", Retrospective: retrospective}
		if err := testDB.Create(&card); err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/retrospectives/{retrospectiveID}/cards/{cardID}/votes", nil)
		rctx.URLParams.Add("retrospectiveID", "invalid")
		rctx.URLParams.Add("cardID", card.ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		env.AddCardVote(rr, r)

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

	t.Run("with wrong retrospectiveID", func(t *testing.T) {
		testDB.CleanDatabase()
		env = Env{
			DB: testDB,
		}
		team := models.Team{Name: "Team Bacon"}
		retrospective := models.Retrospective{Name: "Retrospective Bacon", Team: team}
		card := models.Card{Content: "Card Bacon", Type: "negative", Retrospective: retrospective}
		uuid, _ := uuid.NewV4()
		if err := testDB.Create(&card); err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/retrospectives/{retrospectiveID}/cards/{cardID}/votes", nil)
		rctx.URLParams.Add("retrospectiveID", uuid.String())
		rctx.URLParams.Add("cardID", card.ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		env.AddCardVote(rr, r)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		if res.Errors["card_id"] != "card does not belong to given retrospective" {
			t.Fatalf("card_id error should be card does not belong to given retrospective, but got %v", res.Errors["card_id"])
		}
	})

	t.Run("with db error", func(t *testing.T) {
		dbMock.Error = errors.New("Error trying to save card")
		env = Env{
			DB: dbMock,
		}
		team := models.Team{Name: "Team Bacon"}
		retrospective := models.Retrospective{Name: "Retrospective Bacon", Team: team}
		card := models.Card{Content: "Card Bacon", Type: "negative", Retrospective: retrospective}
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/retrospectives/{retrospectiveID}/cards/{cardID}/votes", nil)
		rctx.URLParams.Add("retrospectiveID", retrospective.ID.String())
		rctx.URLParams.Add("cardID", card.ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		env.AddCardVote(rr, r)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if msg := "error trying to save card"; res.Errors["db"] != msg {
			t.Fatalf("handler should return db error: %s, but got %v", msg, res.Errors["db"])
		}
	})

	t.Run("with valid card", func(t *testing.T) {
		testDB.CleanDatabase()
		env = Env{
			DB: testDB,
		}
		team := models.Team{Name: "Team Bacon"}
		retrospective := models.Retrospective{Name: "Retrospective Bacon", Team: team}
		card := models.Card{Content: "Card Bacon", Type: "negative", Retrospective: retrospective}
		if err := testDB.Create(&card); err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/retrospectives/{retrospectiveID}/cards/{cardID}/votes", nil)
		rctx.URLParams.Add("retrospectiveID", card.RetrospectiveID.String())
		rctx.URLParams.Add("cardID", card.ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		res := response{}

		env.AddCardVote(rr, r)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if res.Data.Votes != 1 {
			t.Fatalf("card votes should be 1, but got %v", res.Data.Votes)
		}
	})
}
