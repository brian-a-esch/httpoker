package poker

import (
	"fmt"
	"sort"
)

const (
	handSize int32 = 5
)

type Hand interface {
	rank() int32
}

const (
	highCardRank int32 = iota
	pairRank
	twoPairRank
	threeKindRank
	straightRank
	flushRank
	fullHouseRank
	fourKindRank
	straightFlushRank
)

type StraightFlush struct {
	sortedCards []Card
}

func (s StraightFlush) rank() int32 {
	return straightFlushRank
}

type FourKind struct {
	fourPair []Card
	kicker   Card
}

func (f FourKind) rank() int32 {
	return fourKindRank
}

type FullHouse struct {
	threePair []Card
	twoPair   []Card
}

func (f FullHouse) rank() int32 {
	return fullHouseRank
}

type Flush struct {
	sortedCards []Card
}

func (f Flush) rank() int32 {
	return flushRank
}

type Straight struct {
	sortedCards []Card
}

func (s Straight) rank() int32 {
	return straightRank
}

type ThreeKind struct {
	threePair []Card
	kickers   []Card
}

func (t ThreeKind) rank() int32 {
	return threeKindRank
}

type TwoPair struct {
	highPair []Card
	lowPair  []Card
	kicker   Card
}

func (p TwoPair) rank() int32 {
	return twoPairRank
}

type Pair struct {
	pair    []Card
	kickers []Card
}

func (p Pair) rank() int32 {
	return pairRank
}

type HighCard struct {
	// TODO this needs to be sorted cards
	sortedCards []Card
}

func (h HighCard) rank() int32 {
	return highCardRank
}

func solveForStraightFlushOrFlush(cards []Card) (*StraightFlush, *Flush) {
	suitsMap := make(map[Suit][]Card)
	for _, card := range cards {
		suitsMap[card.suit] = append(suitsMap[card.suit], card)
	}

	// Here we assume that there is only one flush in the hand
	for _, suitedCards := range suitsMap {
		if len(suitedCards) >= 5 {
			straightMaybe := solveForStraight(suitedCards)
			if straightMaybe != nil {
				return &StraightFlush{sortedCards: straightMaybe.sortedCards}, nil
			}

			sort.Slice(suitedCards, func(i, j int) bool {
				return suitedCards[i].value > suitedCards[j].value
			})
			return nil, &Flush{sortedCards: suitedCards[:5]}
		}
	}

	return nil, nil
}

func solveForStraight(cards []Card) *Straight {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].value > cards[j].value
	})

	result := make([]Card, 0, 5)
	for _, card := range cards {
		if len(result) == 0 || result[len(result)-1].value-1 == card.value {
			result = append(result, card)
		} else if result[len(result)-1].value == card.value {
			continue
		} else {
			result = result[:0]
			result = append(result, card)
		}
	}

	if cards[0].value == Ace && result[len(result)-1].value == Two {
		result = append(result, cards[0])
	}

	if len(result) >= 5 {
		return &Straight{sortedCards: result[:5]}
	}

	return nil
}

func findAllPairs(cards []Card) ([]Card, []Card, []Card) {
	high := make([]Card, 0)
	low := make([]Card, 0)
	kickers := make([]Card, 0)

	cardsByValue := make([][]Card, NumCardValues)
	for _, card := range cards {
		cardsByValue[card.value-2] = append(cardsByValue[card.value-2], card)
	}

	for i := NumCardValues - 1; i >= 0; i-- {
		if len(cardsByValue[i]) > 1 && len(high) == 0 {
			high = cardsByValue[i]
		} else if len(cardsByValue[i]) > 1 && len(low) == 0 {
			low = cardsByValue[i]
		} else {
			kickers = append(kickers, cardsByValue[i]...)
		}
	}

	if len(high) < len(low) {
		return low, high, kickers
	}

	return high, low, kickers
}

// SolveHand takes cards and builds the appropriate hand. Needs to be called with [5, 7] cards
// If called with incorrect number of cards an error is returned
func SolveHand(cards []Card) (Hand, error) {
	if len(cards) < 5 || len(cards) > 7 {
		return nil, fmt.Errorf("Only support hands with with size in [5, 7], but got %d cards", len(cards))
	}

	highPairs, lowPairs, kickers := findAllPairs(cards)
	straightFlush, flush := solveForStraightFlushOrFlush(cards)

	if straightFlush != nil {
		return *straightFlush, nil
	}

	if len(highPairs) == 4 {
		if len(lowPairs) == 0 {
			return FourKind{fourPair: highPairs, kicker: kickers[0]}, nil
		} else if len(kickers) == 0 {
			return FourKind{fourPair: highPairs, kicker: lowPairs[0]}, nil
		} else if lowPairs[0].value < kickers[0].value {
			return FourKind{fourPair: highPairs, kicker: kickers[0]}, nil
		} else {
			return FourKind{fourPair: highPairs, kicker: lowPairs[0]}, nil
		}
	}

	if len(highPairs) == 3 && len(lowPairs) == 2 {
		return FullHouse{threePair: highPairs, twoPair: lowPairs}, nil
	}

	if flush != nil {
		return *flush, nil
	}

	straight := solveForStraight(cards)
	if straight != nil {
		return *straight, nil
	}

	if len(highPairs) == 3 {
		return ThreeKind{threePair: highPairs, kickers: kickers[:2]}, nil
	}

	if len(highPairs) == 2 && len(lowPairs) == 2 {
		return TwoPair{highPair: highPairs, lowPair: lowPairs, kicker: kickers[0]}, nil
	}

	if len(highPairs) == 2 {
		return Pair{pair: highPairs, kickers: kickers[:3]}, nil
	}

	return HighCard{sortedCards: kickers[:5]}, nil
}

func kickerCompare(lhsKickers []Card, rhsKickers []Card) int {
	if len(lhsKickers) != len(rhsKickers) {
		panic("Cannot compare kickers when lengths are different")
	}

	for i, lhs := range lhsKickers {
		if lhs.value != rhsKickers[i].value {
			return int(rhsKickers[i].value - lhs.value)
		}
	}

	return 0
}

// CompareHands takes two properly constructed hands and determines who wins.
// Returns a negative number is lhs wins, positive if rhs win, and 0 if they chop
func CompareHands(lhsHand Hand, rhsHand Hand) int {
	if lhsHand.rank() < rhsHand.rank() {
		return 1
	} else if lhsHand.rank() > rhsHand.rank() {
		return -1
	}

	switch lhs := lhsHand.(type) {
	case StraightFlush:
		rhs := rhsHand.(StraightFlush)
		return int(rhs.sortedCards[0].value - lhs.sortedCards[0].value)
	case FourKind:
		rhs := rhsHand.(FourKind)
		if lhs.fourPair[0].value != rhs.fourPair[0].value {
			return int(rhs.fourPair[0].value - lhs.fourPair[0].value)
		}
		return int(rhs.kicker.value - lhs.kicker.value)
	case FullHouse:
		rhs := rhsHand.(FullHouse)
		if lhs.threePair[0].value != rhs.threePair[0].value {
			return int(rhs.threePair[0].value - lhs.threePair[0].value)
		}
		return int(rhs.twoPair[0].value - lhs.twoPair[0].value)
	case Flush:
		rhs := rhsHand.(Flush)
		return int(rhs.sortedCards[0].value - lhs.sortedCards[0].value)
	case Straight:
		rhs := rhsHand.(Straight)
		return int(rhs.sortedCards[0].value - lhs.sortedCards[0].value)
	case ThreeKind:
		rhs := rhsHand.(ThreeKind)
		if lhs.threePair[0].value != rhs.threePair[0].value {
			return int(rhs.threePair[0].value - lhs.threePair[0].value)
		}
		return kickerCompare(lhs.kickers, rhs.kickers)
	case TwoPair:
		rhs := rhsHand.(TwoPair)
		if lhs.highPair[0].value != rhs.highPair[0].value {
			return int(rhs.highPair[0].value - lhs.highPair[0].value)
		}
		if lhs.lowPair[0].value != rhs.lowPair[0].value {
			return int(rhs.lowPair[0].value - lhs.lowPair[0].value)
		}
		return int(rhs.kicker.value - lhs.kicker.value)
	case Pair:
		rhs := rhsHand.(Pair)
		if lhs.pair[0].value != rhs.pair[0].value {
			return int(rhs.pair[0].value - lhs.pair[0].value)
		}
		return kickerCompare(lhs.kickers, rhs.kickers)
	case HighCard:
		rhs := rhsHand.(HighCard)
		return kickerCompare(lhs.sortedCards, rhs.sortedCards)
	default:
		panic("Got a Hand which did not match any expected types")
	}
}
