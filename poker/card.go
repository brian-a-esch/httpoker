package poker

// Suit is the suit of a card
type Suit int

const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
)

func (s Suit) String() string {
	return [...]string{
		"Spades",
		"Hearts",
		"Diamonds",
		"Clubs",
		"Foo",
	}[s]
}

// CardValue is the value of a card, i.e. 2, Ace, King etc. We make them comparable
// by performing less than on them
type CardValue int

const (
	Two CardValue = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

func (v CardValue) String() string {
	return [...]string{
		"Two",
		"Three",
		"Four",
		"Five",
		"Six",
		"Seven",
		"Eight",
		"Nine",
		"Ten",
		"Jack",
		"Queen",
		"King",
		"Ace",
	}[v-2]
}

// Card is a struct with a suit and a value
type Card struct {
	suit  Suit
	value CardValue
}

const (
	// NumCardValues is the number of possible card values
	NumCardValues = 13
	// NumSuits is the umber of possible card suits
	NumSuits = 4
)
