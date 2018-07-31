package models

import (
	"github.com/gobuffalo/uuid"
	"github.com/jinzhu/gorm"
)

type Team struct {
	Model
	Name string `json:"name" db:"name"`
}

func FetchTeamByID(team_id uuid.UUID, team *Team, db *gorm.DB) error {
	return db.Where("id = ?", team_id).First(team).Error
}
