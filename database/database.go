package database

import (
	"github.com/gobuffalo/uuid"
	"github.com/nicolaszein/go-retro/models"
)

type Database interface {
	Create(interface{}) error
	CleanDatabase()
	FetchTeams(interface{}) error
	FetchTeamByID(uuid.UUID, *models.Team) error
	FetchRetrospectiveByID(uuid.UUID, *models.Retrospective) error
}
