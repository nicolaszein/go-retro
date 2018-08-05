package database

import (
	"github.com/jinzhu/gorm"
)

type Database interface {
	Create(interface{}) error
}

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

type Mock struct {
	Error error
}

func (m Mock) Create(interface{}) error {
	return m.Error
}
