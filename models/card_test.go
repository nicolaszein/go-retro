package models

import (
	"testing"
)

func TestCardAddVote(t *testing.T) {
	card := Card{}

	card.AddVote()

	if card.Votes != 1 {
		t.Fatalf("card votes should be 1, but got %v", card.Votes)
	}
}
