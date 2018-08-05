package database

import (
	"github.com/jinzhu/gorm"
	"github.com/nicolaszein/go-retro/models"
)

// Postgres implementing Database interface
type Postgres struct {
	DB *gorm.DB
}

func (p Postgres) Create(s interface{}) error {
	return p.DB.Create(s).Error
}

func NewPostgres(url string) (*Postgres, error) {
	db, err := gorm.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &Postgres{db}, nil
}

func (p Postgres) CleanDatabase() {
	p.DB.Unscoped().Delete(&models.Team{})
}
