package poker

import (
	"sort"
)

const (
	handSize int32 = 5
)

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

type Hand interface {
	rank() int32
}

type StraightFlush struct {
	sortedCards []Card
}

func (s StraightFlush) rank() int32 {
	return straightFlushRank
}

/*
func solveForStraightFlush(cards []Card) (StraightFlush, bool) {
	Cant just use flush and straight solvers, consider
	2H 3H 4H 5H 6H 4C JH
	there are two straight solutions and two flush solutions. The solve methods
	just find the best straight/flush, they don't list all solutions
}
*/

type Flush struct {
	sortedCards []Card
}

func (f Flush) rank() int32 {
	return flushRank
}

func solveForFlush(cards []Card) *Flush {
	suitsMap := make(map[Suit][]Card)
	for _, card := range cards {
		suitsMap[card.suit] = append(suitsMap[card.suit], card)
	}

	// Here we assume that there is only one flush in the hand
	for _, cards := range suitsMap {
		if len(cards) >= 5 {
			sort.Slice(cards, func(i, j int) bool {
				return cards[i].value > cards[j].value
			})

			return &Flush{sortedCards: cards[:5]}
		}
	}

	return nil
}

type Straight struct {
	sortedCards []Card
}

func (s Straight) rank() int32 {
	return straightRank
}

func solveForStraight(cards []Card) *Straight {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].value > cards[i].value
	})

	/*
		result := make([]Card, 5)
		for i, card := range cards {
			if i == 0 || lastCard.value {
				streakStart = i
				lastCard = &cards
			}
		}
	*/
	return nil
}

type HighCard struct {
	// TODO this needs to be sorted cards
	highCard Card
}

func (h HighCard) rank() int32 {
	return highCardRank
}

func solveForHighCard(cards []Card) *HighCard {
	var high Card
	for i, card := range cards {
		if i == 0 || card.value > high.value {
			high = card
		}
	}

	return &HighCard{highCard: high}
}

func SolveHand(cards []Card) Hand {
	if len(cards) == 0 {
		panic("Do not support empty hands")
	} else if len(cards) > 7 {
		panic("Do not support hands larger than 7. Assumptions made in solvers that hand size does not exceed 7")
	}

	flush := solveForFlush(cards)
	if flush != nil {
		return flush
	}

	straight := solveForStraight(cards)
	if flush != nil {
		return straight
	}

	return solveForHighCard(cards)
}
