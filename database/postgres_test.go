package database

import (
	"testing"

	"github.com/gobuffalo/uuid"
	"github.com/nicolaszein/go-retro/models"
)

func TestFetchTeamByID(t *testing.T) {
	t.Run("with a persisted team", func(t *testing.T) {
		testDB.CleanDatabase()
		persistedTeam := models.Team{Name: "Team Bacon"}

		if err := testDB.Create(&persistedTeam); err != nil {
			t.Fatalf("Failed creating persistedTeam %v", err)
		}

		team := models.Team{}
		if err := testDB.FetchTeamByID(persistedTeam.ID, &team); err != nil {
			t.Fatalf("err should be nil, but got %v", err)
		}
	})

	t.Run("with no team", func(t *testing.T) {
		testDB.CleanDatabase()
		team := models.Team{}
		teamID, err := uuid.NewV4()
		if err != nil {
			t.Fatalf("error trying to create uuid")
		}

		if err := testDB.FetchTeamByID(teamID, &team); err == nil {
			t.Fatalf("team should be nil, but got %v", team)
		}
	})
}

func TestFetchRetrospectivesByTeamID(t *testing.T) {
	t.Run("with restrospectives with same team", func(t *testing.T) {
		testDB.CleanDatabase()
		team := models.Team{Name: "Team Bacon"}
		if err := testDB.Create(&team); err != nil {
			t.Fatalf("Failed creating persistedTeam")
		}
		retrospective := models.Retrospective{Name: "Retrospective 1", TeamID: team.ID}
		if err := testDB.Create(&retrospective); err != nil {
			t.Fatalf("Failed creating retrospective %v", err)
		}
		retrospective2 := models.Retrospective{Name: "Retrospective 2", TeamID: team.ID}
		if err := testDB.Create(&retrospective2); err != nil {
			t.Fatalf("Failed creating retrospective %v", err)
		}

		retrospectives := []models.Retrospective{}
		if err := testDB.FetchRestrospectivesByTeamID(team.ID, &retrospectives); err != nil {
			t.Fatalf("err should be nil, but got %v", err)
		}

		if len(retrospectives) != 2 {
			t.Fatalf("retrospectives len should be 2, but got %v", len(retrospectives))
		}
	})

	t.Run("with retrospectives with different team", func(t *testing.T) {
		testDB.CleanDatabase()
		team := models.Team{Name: "Team Bacon"}
		if err := testDB.Create(&team); err != nil {
			t.Fatalf("Failed creating persistedTeam")
		}
		team2 := models.Team{Name: "Team Fitness"}
		if err := testDB.Create(&team2); err != nil {
			t.Fatalf("Failed creating persistedTeam")
		}
		retrospective := models.Retrospective{Name: "Retrospective 1", TeamID: team.ID}
		if err := testDB.Create(&retrospective); err != nil {
			t.Fatalf("Failed creating retrospective %v", err)
		}
		retrospective2 := models.Retrospective{Name: "Retrospective 2", TeamID: team2.ID}
		if err := testDB.Create(&retrospective2); err != nil {
			t.Fatalf("Failed creating retrospective %v", err)
		}

		retrospectives := []models.Retrospective{}
		if err := testDB.FetchRestrospectivesByTeamID(team.ID, &retrospectives); err != nil {
			t.Fatalf("err should be nil, but got %v", err)
		}

		if len(retrospectives) != 1 {
			t.Fatalf("retrospectives len should be 1, but got %v", len(retrospectives))
		}
	})

	t.Run("with no retrospectives for team", func(t *testing.T) {
		testDB.CleanDatabase()
		team := models.Team{Name: "Team Bacon"}
		if err := testDB.Create(&team); err != nil {
			t.Fatalf("Failed creating persistedTeam")
		}

		retrospectives := []models.Retrospective{}
		if err := testDB.FetchRestrospectivesByTeamID(team.ID, &retrospectives); err != nil {
			t.Fatalf("err should be nil, but got %v", err)
		}

		if len(retrospectives) != 0 {
			t.Fatalf("retrospectives len should be 0, but got %v", len(retrospectives))
		}
	})
}
