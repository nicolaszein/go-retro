package models

type Team struct {
	Model
	Name string `json:"name" db:"name" valid:"required"`
}
