package poker

import "testing"

func TestDeckDeal(t *testing.T) {
	deck := NewDeck()
	cards := make(map[Card]int)
	for i := 0; i < NumCards; i++ {
		c, err := deck.DealCard()
		if err != nil {
			t.Error()
		}
		cards[c]++
	}

	for _, count := range cards {
		if count != 1 {
			t.Error()
		}
	}
	if len(cards) != NumCards {
		t.Error()
	}

	_, err := deck.DealCard()
	if err == nil {
		t.Error()
	}
}

func TestDealAndShuffle(t *testing.T) {
	deck := NewDeck()
	card, _ := deck.DealCard()
	deck.Shuffle()
	foundCard := false

	for deck.Len() > 0 {
		c, err := deck.DealCard()
		if err != nil {
			t.Error()
		}

		if card == c {
			foundCard = true
		}
	}

	if !foundCard {
		t.Error()
	}
}
