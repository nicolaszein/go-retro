package models

import (
	"github.com/gobuffalo/uuid"
	"github.com/jinzhu/gorm"
)

type Retrospective struct {
	Model
	Name   string    `json:"name" db:"name" valid:"required"`
	TeamID uuid.UUID `json:"team_id" db:"team_id"`
	Team   Team      `json:"team"`
}

func FetchRestrospectivesByTeamID(team_id uuid.UUID, retrospectives *[]Retrospective, db *gorm.DB) error {
	return db.Where("team_id = ?", team_id).Find(retrospectives).Error
}
