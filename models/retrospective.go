package models

import (
	"github.com/gobuffalo/uuid"
)

type Retrospective struct {
	Model
	Name   string    `json:"name" db:"name" valid:"required"`
	TeamID uuid.UUID `json:"team_id" db:"team_id"`
	Team   Team      `json:"team"`
}
