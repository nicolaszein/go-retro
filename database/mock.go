package database

import (
	"github.com/gobuffalo/uuid"
	"github.com/nicolaszein/go-retro/models"
)

type Mock struct {
	Error              error
	FetchTeamByIDError error
}

func (m Mock) Create(interface{}) error {
	return m.Error
}

func (m Mock) CleanDatabase() {}

// Teams
func (m Mock) FetchTeams(interface{}) error {
	return m.Error
}

func (m Mock) FetchTeamByID(team_id uuid.UUID, team *models.Team) error {
	return m.FetchTeamByIDError
}

// Retrospectives
func (m Mock) FetchRestrospectivesByTeamID(team_id uuid.UUID, retrospectives *[]models.Retrospective) error {
	return m.Error
}
