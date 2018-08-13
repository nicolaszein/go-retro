package database

import (
	"github.com/gobuffalo/uuid"
	"github.com/nicolaszein/go-retro/models"
)

type Mock struct {
	Error                       error
	FetchTeamByIDError          error
	FetchRetrospectiveByIDError error
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

func (m Mock) FetchRetrospectiveByID(retrospective_id uuid.UUID, retrospective *models.Retrospective) error {
	return m.FetchRetrospectiveByIDError
}

// Cards
func (m Mock) FetchCardByID(card_id uuid.UUID, card *models.Card) error {
	return m.Error
}
