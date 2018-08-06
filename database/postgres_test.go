package database

import (
	"os"
	"testing"

	"github.com/gobuffalo/uuid"
	"github.com/nicolaszein/go-retro/models"
)

func TestCreate(t *testing.T) {
	t.Run("with valid interface", func(t *testing.T) {
		testDB.CleanDatabase()
		team := models.Team{Name: "Team Bacon"}

		if err := testDB.Create(&team); err != nil {
			t.Fatalf("Failed creating team %v", err)
		}
	})

	t.Run("with invalid interface", func(t *testing.T) {
		testDB.CleanDatabase()
		retrospective := models.Retrospective{}

		if err := testDB.Create(&retrospective); err == nil {
			t.Fatalf("err should have value, but got nil")
		}
	})
}

func TestNewConnection(t *testing.T) {
	t.Run("with valid db", func(t *testing.T) {
		url := os.Getenv("TEST_DATABASE_URL")

		_, err := NewPostgres(url)
		if err != nil {
			t.Fatalf("err should be nil, but got %v", err)
		}
	})

	t.Run("with invalid db", func(t *testing.T) {
		url := "invalid_database"

		db, _ := NewPostgres(url)
		if db != nil {
			t.Fatalf("db should be nil, but got %v", db)
		}
	})
}

func TestCleanDatabase(t *testing.T) {
	t.Run("clean teams", func(t *testing.T) {
		persistedTeam := models.Team{Name: "Team Bacon"}
		if err := testDB.Create(&persistedTeam); err != nil {
			t.Fatalf("Failed creating persistedTeam %v", err)
		}

		testDB.CleanDatabase()

		teams := []models.Team{}
		if err := testDB.DB.Find(&teams).Error; err != nil {
			t.Fatalf("Error trying to get all teams! Error %v", err)
		}
		if len(teams) > 0 {
			t.Fatalf("teams len should be 0, but got %v", len(teams))
		}
	})

	t.Run("clean retrospectives", func(t *testing.T) {
		persistedTeam := models.Team{Name: "Team Bacon"}
		if err := testDB.Create(&persistedTeam); err != nil {
			t.Fatalf("Failed creating persistedTeam %v", err)
		}
		persistedRetrospective := models.Retrospective{Name: "Retrospective Bacon", TeamID: persistedTeam.ID}
		if err := testDB.Create(&persistedRetrospective); err != nil {
			t.Fatalf("Failed creating persistedRetrospective %v", err)
		}

		testDB.CleanDatabase()

		retrospectives := []models.Retrospective{}
		if err := testDB.DB.Find(&retrospectives).Error; err != nil {
			t.Fatalf("Error trying to get all retrospectives! Error %v", err)
		}
		if len(retrospectives) > 0 {
			t.Fatalf("retrospectives len should be 0, but got %v", len(retrospectives))
		}
	})
}

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
