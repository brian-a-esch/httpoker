package poker

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

type Card struct {
	suit  Suit
	value CardValue
}

func CardCompare(lhs Card, rhs Card) int32 {
	if lhs.value < rhs.value {
		return -1
	} else if lhs.value > rhs.value {
		return 1
	} else {
		return 0
	}
}
