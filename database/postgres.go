package database

import (
	"github.com/gobuffalo/uuid"
	"github.com/jinzhu/gorm"
	"github.com/nicolaszein/go-retro/models"
)

// Postgres implementing Database interface
type Postgres struct {
	DB *gorm.DB
}

func NewPostgres(url string) (*Postgres, error) {
	db, err := gorm.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &Postgres{db}, nil
}

func (p Postgres) Create(s interface{}) error {
	return p.DB.Create(s).Error
}

func (p Postgres) CleanDatabase() {
	p.DB.Unscoped().Delete(&models.Card{})
	p.DB.Unscoped().Delete(&models.Retrospective{})
	p.DB.Unscoped().Delete(&models.Team{})
}

// Teams
func (p Postgres) FetchTeams(s interface{}) error {
	return p.DB.Order("name").Find(s).Error
}

func (p Postgres) FetchTeamByID(team_id uuid.UUID, team *models.Team) error {
	return p.DB.Where("id = ?", team_id).First(team).Error
}

// Retrospectives
func (p Postgres) FetchRestrospectivesByTeamID(team_id uuid.UUID, retrospectives *[]models.Retrospective) error {
	return p.DB.Where("team_id = ?", team_id).Find(retrospectives).Error
}

func (p Postgres) FetchRetrospectiveByID(retrospective_id uuid.UUID, retrospective *models.Retrospective) error {
	return p.DB.Where("id = ?", retrospective_id).First(retrospective).Error
}
