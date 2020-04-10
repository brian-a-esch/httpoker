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

func TestCardCompare(t *testing.T) {
	c1 := Card{suit: Spades, value: Two}
	c2 := Card{suit: Hearts, value: Two}
	c3 := Card{suit: Hearts, value: Three}

	if CardCompare(c1, c2) != 0 {
		t.Error()
	}

	if CardCompare(c1, c3) >= 0 {
		t.Error()
	}

	if CardCompare(c3, c1) <= 0 {
		t.Error()
	}
}
