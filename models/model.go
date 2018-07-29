package models

import (
	"time"

	"github.com/gobuffalo/uuid"
	"github.com/jinzhu/gorm"
)

type Model struct {
	ID        *uuid.UUID `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (model *Model) BeforeCreate(scope *gorm.Scope) error {
	id, _ := uuid.NewV4()
	scope.SetColumn("ID", &id)
	model.ID = &id
	return nil
}
