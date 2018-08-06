package database

import (
	"errors"
	"testing"

	"github.com/gobuffalo/uuid"
	"github.com/nicolaszein/go-retro/models"
)

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
