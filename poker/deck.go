package poker

import (
	"fmt"
	"math/rand"
	"time"
)

// NumCards in the deck
const NumCards int = 52

// Deck keeps track of cards
type Deck struct {
	cards  []Card
	dealt  []Card
	random *rand.Rand
}

// NewDeck creates a new Deck
func NewDeck() Deck {
	cards := make([]Card, 0, NumCards)

	var suits = [...]Suit{
		Spades,
		Hearts,
		Diamonds,
		Clubs,
	}

	var values = [...]CardValue{
		Two,
		Three,
		Four,
		Five,
		Six,
		Seven,
		Eight,
		Nine,
		Ten,
		Jack,
		Queen,
		King,
		Ace,
	}

	for _, suit := range suits {
		for _, value := range values {
			cards = append(cards, Card{suit: suit, value: value})
		}
	}

	if len(cards) != NumCards {
		panic("Did not construct correct number of cards")
	}

	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	deck := Deck{cards: cards, dealt: make([]Card, 0), random: random}
	deck.Shuffle()
	return deck
}

// Shuffle takes all the dealt and un-dealt cards and s
func (deck *Deck) Shuffle() {
	deck.cards = append(deck.cards, deck.dealt...)
	if len(deck.cards) != NumCards {
		panic("Deck somehow does not have correct nubmer of cards")
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck.cards), func(i, j int) { deck.cards[i], deck.cards[j] = deck.cards[j], deck.cards[i] })
}

// DealCard deals a card and returns an error if there are no more cards to deal
func (deck *Deck) DealCard() (Card, error) {
	if len(deck.cards) <= 0 {
		return Card{}, fmt.Errorf("Deck is out of cards")
	}

	var card Card
	card, deck.cards = deck.cards[len(deck.cards)-1], deck.cards[:len(deck.cards)-1]
	deck.dealt = append(deck.dealt, card)
	return card, nil
}

// Len returns the number of cards remaining in the deck
func (deck *Deck) Len() int {
	return len(deck.cards)
}
