package models

import (
	"testing"
)

func TestFetchByTeamID(t *testing.T) {
	t.Run("with persisteds restrospectives", func(t *testing.T) {
		cleanDatabase(testDB)
		team := Team{Name: "Team Bacon"}

		if err := testDB.Create(&team).Error; err != nil {
			t.Fatalf("Failed creating persistedTeam %v", err)
		}

		retrospective := Retrospective{Name: "Retrospective 1", TeamID: team.ID}
		if err := testDB.Create(&retrospective).Error; err != nil {
			t.Fatalf("Failed creating retrospective %v", err)
		}

		retrospective2 := Retrospective{Name: "Retrospective 2", TeamID: team.ID}
		if err := testDB.Create(&retrospective2).Error; err != nil {
			t.Fatalf("Failed creating retrospective %v", err)
		}

		retrospectives := []Retrospective{}
		if err := FetchRestrospectivesByTeamID(team.ID, &retrospectives, testDB); err != nil {
			t.Fatalf("err should be nil, but got %v", err)
		}

		if len(retrospectives) != 2 {
			t.Fatalf("retrospectives len should be 2, but got %v", len(retrospectives))
		}
	})
}
