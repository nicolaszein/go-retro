package database

import (
	"errors"
	"testing"

	"github.com/gobuffalo/uuid"
	"github.com/nicolaszein/go-retro/models"
)

func TestMockCreate(t *testing.T) {
	dbMock.CleanDatabase()
	team := models.Team{}
	dbMock.Error = errors.New("error trying to create team")

	if err := dbMock.Create(&team).Error; err == nil {
		t.Fatalf("err should have value but got nil")
	}
}

func TestMockFetchTeams(t *testing.T) {
	dbMock.CleanDatabase()
	dbMock.Error = errors.New("Error trying to fetch teams")

	teams := []models.Team{}
	if err := dbMock.FetchTeams(&teams); err == nil {
		t.Fatalf("teams should be nil, but got %v", teams)
	}
}

func TestMockFetchTeamByID(t *testing.T) {
	dbMock.CleanDatabase()
	dbMock.FetchTeamByIDError = errors.New("Error trying to fetch teams")
	uuid, _ := uuid.NewV4()

	team := models.Team{}
	if err := dbMock.FetchTeamByID(uuid, &team); err == nil {
		t.Fatalf("teams should be nil, but got %v", team)
	}
}

func TestMockFetchRetrospectivesByTeamID(t *testing.T) {
	t.Run("with restrospectives with same team", func(t *testing.T) {
		dbMock.CleanDatabase()
		dbMock.Error = errors.New("Error trying to fetch retrospectives")
		teamID, err := uuid.NewV4()
		if err != nil {
			t.Fatalf("error trying to create uuid")
		}

		retrospectives := []models.Retrospective{}
		if err := dbMock.FetchRestrospectivesByTeamID(teamID, &retrospectives); err == nil {
			t.Fatalf("retrospectives should be nil, but got %v", retrospectives)
		}
	})
}

func TestMockFetchRetrospectiveByID(t *testing.T) {
	dbMock.CleanDatabase()
	dbMock.FetchRetrospectiveByIDError = errors.New("Error trying to fetch retrospective")
	retrospectiveID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("error trying to create uuid")
	}

	retrospective := models.Retrospective{}
	if err := dbMock.FetchRetrospectiveByID(retrospectiveID, &retrospective); err == nil {
		t.Fatalf("retrospective should be nil, but got %v", retrospective)
	}
}
