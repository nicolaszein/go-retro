package models

import (
	"testing"
)

func TestFetchByID(t *testing.T) {
	t.Run("with a persisted team", func(t *testing.T) {
		cleanDatabase(testDB)
		persisted_team := Team{Name: "Team Bacon"}
		testDB.Create(&persisted_team)

		_, ok := FetchTeamByID(persisted_team.ID, testDB)

		if !ok {
			t.Errorf(`team shouldn't be nil`)
		}
	})
}
