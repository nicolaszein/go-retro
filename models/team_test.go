package models

import (
	"testing"
)

func TestFetchByID(t *testing.T) {
	t.Run("with a persisted team", func(t *testing.T) {
		cleanDatabase(testDB)
		persistedTeam := Team{Name: "Team Bacon"}

		if err := testDB.Create(&persistedTeam).Error; err != nil {
			t.Fatalf("Failed creating persistedTeam %v", err)
		}

		team := Team{}
		if err := FetchTeamByID(persistedTeam.ID, &team, testDB); err != nil {
			t.Fatalf("err should be nil, but got %v", err)
		}
	})
}
