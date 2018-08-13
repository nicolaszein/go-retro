package models

import "github.com/gobuffalo/uuid"

type Card struct {
	Model
	Content         string        `json:"content" db:"content" valid:"required"`
	RetrospectiveID uuid.UUID     `json:"retrospective_id" db:"retrospective_id"`
	Retrospective   Retrospective `json:"retrospective" valid:"-"`
	Votes           int           `json:"votes" db:"votes"`
	Type            string        `json:"type" db:"type" valid:"matches(^positive$|^negative$),required"`
}

func (card *Card) AddVote() {
	card.Votes += 1
}
