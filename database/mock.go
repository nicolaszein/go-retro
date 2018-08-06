package database

import (
	"github.com/gobuffalo/uuid"
	"github.com/nicolaszein/go-retro/models"
)

type Mock struct {
	Error error
}

func (m Mock) Create(interface{}) error {
	return m.Error
}
func (m Mock) CleanDatabase() {}

func (m Mock) FetchRestrospectivesByTeamID(team_id uuid.UUID, retrospectives *[]models.Retrospective) error {
	return m.Error
}
