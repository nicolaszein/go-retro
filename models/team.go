package models

import (
	"github.com/gobuffalo/uuid"
	"github.com/jinzhu/gorm"
)

type Team struct {
	Model
	Name string `json:"name" db:"name"`
}

func FetchTeamByID(team_id *uuid.UUID, db *gorm.DB) (team Team, ok bool) {
	team = Team{}
	db.Where("id = ?", team_id).First(&team)
	ok = team.ID.String() == team_id.String()
	return team, ok
}
