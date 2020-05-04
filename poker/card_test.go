package poker

import "testing"

func TestSuitString(t *testing.T) {
	someSuit := Spades
	if someSuit.String() != "Spades" {
		t.Error()
	}
}

func TestCardString(t *testing.T) {
	c := Three
	if c.String() != "Three" {
		t.Error()
	}
}
