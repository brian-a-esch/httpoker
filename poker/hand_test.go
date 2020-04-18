package poker

import (
	"reflect"
	"testing"
)

func buildCards(cards []string) []Card {
	result := make([]Card, 0)
	for _, card := range cards {
		if len(card) != 2 {
			panic("Only handle two char card shortcut")
		}

		var value CardValue
		switch card[0] {
		case '2':
			value = Two
		case '3':
			value = Three
		case '4':
			value = Four
		case '5':
			value = Five
		case '6':
			value = Six
		case '7':
			value = Seven
		case '8':
			value = Eight
		case '9':
			value = Nine
		case 'T':
			value = Ten
		case 'J':
			value = Jack
		case 'Q':
			value = Queen
		case 'K':
			value = King
		case 'A':
			value = Ace
		default:
			panic("Unknown card value")
		}

		var suit Suit
		switch card[1] {
		case 'S':
			suit = Spades
		case 'H':
			suit = Hearts
		case 'D':
			suit = Diamonds
		case 'C':
			suit = Clubs
		default:
			panic("Unknown suit")
		}

		result = append(result, Card{suit: suit, value: value})
	}

	return result
}

func cardsMakePair(cards []Card, value CardValue, expectedLength int) bool {
	if len(cards) != expectedLength {
		return false
	}

	for _, c := range cards {
		if c.value != value {
			return false
		}
	}

	return true
}

func TestHighCardBasic(t *testing.T) {
	cards := buildCards([]string{"2H", "4C", "8D", "TS", "AS", "3H", "9C"})
	hand, _ := SolveHand(cards)
	highCard := hand.(HighCard)

	expected := buildCards([]string{"AS", "TS", "9C", "8D", "4C"})
	if !reflect.DeepEqual(highCard.sortedCards, expected) {
		t.Error()
	}

	if highCard.rank() != highCardRank {
		t.Error()
	}
}

func TestFlushBasic(t *testing.T) {
	cards := buildCards([]string{"AS", "2S", "4S", "8S", "TS"})
	hand, _ := SolveHand(cards)
	flush := hand.(Flush)

	expected := buildCards([]string{"AS", "TS", "8S", "4S", "2S"})

	if !reflect.DeepEqual(expected, flush.sortedCards) {
		t.Error()
	}

	if flush.rank() != flushRank {
		t.Error()
	}
}

func TestFlushLong(t *testing.T) {
	cards := buildCards([]string{"AS", "2S", "4S", "8S", "TS", "3S", "9S"})
	hand, _ := SolveHand(cards)
	flush := hand.(Flush)

	expected := buildCards([]string{"AS", "TS", "9S", "8S", "4S"})
	if !reflect.DeepEqual(expected, flush.sortedCards) {
		t.Error()
	}
}

func TestStraightBasic(t *testing.T) {
	cards := buildCards([]string{"2H", "3C", "6S", "TS", "QC", "4D", "5D"})
	hand, _ := SolveHand(cards)
	straight := hand.(Straight)

	expected := buildCards([]string{"6S", "5D", "4D", "3C", "2H"})
	if !reflect.DeepEqual(expected, straight.sortedCards) {
		t.Error()
	}

	if straight.rank() != straightRank {
		t.Error()
	}
}

func TestStraightDuplicateCards(t *testing.T) {
	cards := buildCards([]string{"2H", "3C", "6S", "TS", "5S", "4D", "5D"})
	hand, _ := SolveHand(cards)
	straight := hand.(Straight)

	possibility1 := buildCards([]string{"6S", "5S", "4D", "3C", "2H"})
	possibility2 := buildCards([]string{"6S", "5D", "4D", "3C", "2H"})

	if !reflect.DeepEqual(straight.sortedCards, possibility1) && !reflect.DeepEqual(straight.sortedCards, possibility2) {
		t.Error()
	}
}

func TestStraightLong(t *testing.T) {
	cards := buildCards([]string{"2H", "3C", "6S", "8S", "4D", "5D", "7H"})
	hand, _ := SolveHand(cards)
	straight := hand.(Straight)

	expected := buildCards([]string{"8S", "7H", "6S", "5D", "4D"})
	if !reflect.DeepEqual(straight.sortedCards, expected) {
		t.Error("Need to select the highest straight possibility")
	}
}

func TestStraightAceLow(t *testing.T) {
	cards := buildCards([]string{"2H", "3C", "AS", "TS", "QC", "4D", "5D"})
	hand, _ := SolveHand(cards)
	straight := hand.(Straight)

	expected := buildCards([]string{"5D", "4D", "3C", "2H", "AS"})
	if !reflect.DeepEqual(straight.sortedCards, expected) {
		t.Error("Ace needs to be able to wrap around for low straight")
	}
}

func TestStraightAceWraparound(t *testing.T) {
	cards := buildCards([]string{"AS", "KS", "QC", "2H", "3D"})
	hand, _ := SolveHand(cards)
	_, castOkay := hand.(Straight)

	if castOkay {
		t.Error("Ace should not wrap around to form straight")
	}
}

func TestFourKind(t *testing.T) {
	cards := buildCards([]string{"AH", "JS", "JC", "JD", "JH", "3S"})
	hand, _ := SolveHand(cards)
	fourKind := hand.(FourKind)

	if fourKind.rank() != fourKindRank {
		t.Error()
	}
	if !cardsMakePair(fourKind.fourPair, Jack, 4) {
		t.Error()
	}
	if fourKind.kicker.value != Ace {
		t.Error("Did not get Ace as kicker")
	}
}

func TestFourKindKickers(t *testing.T) {
	cards := buildCards([]string{"2H", "2C", "2S", "2D", "AH", "AS", "5H"})
	hand, _ := SolveHand(cards)
	fourKind := hand.(FourKind)

	if fourKind.kicker.value != Ace {
		t.Error()
	}

	cards = buildCards([]string{"2H", "2C", "2S", "2D", "AH", "5S", "5H"})
	hand, _ = SolveHand(cards)
	fourKind = hand.(FourKind)
	if fourKind.kicker.value != Ace {
		t.Error()
	}

	cards = buildCards([]string{"2H", "2C", "2S", "2D", "AH", "5S", "6H"})
	hand, _ = SolveHand(cards)
	fourKind = hand.(FourKind)
	if fourKind.kicker.value != Ace {
		t.Error()
	}

	cards = buildCards([]string{"2H", "2C", "2S", "2D", "5H", "5S", "5H"})
	hand, _ = SolveHand(cards)
	fourKind = hand.(FourKind)
	if fourKind.kicker.value != Five {
		t.Error()
	}
}

func TestStraightFlush(t *testing.T) {
	cards := buildCards([]string{"7H", "8H", "9H", "TH", "JH"})
	hand, _ := SolveHand(cards)
	straightFlush := hand.(StraightFlush)

	expected := buildCards([]string{"JH", "TH", "9H", "8H", "7H"})
	if straightFlush.rank() != straightFlushRank {
		t.Error()
	}
	if !reflect.DeepEqual(straightFlush.sortedCards, expected) {
		t.Error()
	}
}

func TestStraightFlushHard(t *testing.T) {
	cards := buildCards([]string{"2H", "3H", "4H", "5H", "6H", "7C", "JH"})
	hand, _ := SolveHand(cards)
	straightFlush := hand.(StraightFlush)

	expected := buildCards([]string{"2H", "3H", "4H", "5H", "6H"})
	if !reflect.DeepEqual(straightFlush.sortedCards, expected) {
		t.Error()
	}
}

func TestFullHouse(t *testing.T) {
	cards := buildCards([]string{"2C", "2H", "2S", "5H", "5S", "AS", "KC"})
	hand, _ := SolveHand(cards)
	fullHouse := hand.(FullHouse)

	if fullHouse.rank() != fullHouseRank {
		t.Error()
	}
	if !cardsMakePair(fullHouse.threePair, Two, 3) {
		t.Error()
	}
	if !cardsMakePair(fullHouse.twoPair, Five, 2) {
		t.Error()
	}
}

func TestThreeKind(t *testing.T) {
	cards := buildCards([]string{"2C", "4H", "5S", "5H", "5S", "AS", "JC"})
	hand, _ := SolveHand(cards)
	threeKind := hand.(ThreeKind)

	expectedKickers := buildCards([]string{"AS", "JC"})
	if threeKind.rank() != threeKindRank {
		t.Error()
	}
	if !cardsMakePair(threeKind.threePair, Five, 3) {
		t.Error()
	}
	if !reflect.DeepEqual(expectedKickers, threeKind.kickers) {
		t.Error()
	}
}

func TestTwoPair(t *testing.T) {
	cards := buildCards([]string{"2C", "4H", "4S", "5H", "5S", "AS", "KC"})
	hand, _ := SolveHand(cards)
	twoPair := hand.(TwoPair)

	if twoPair.rank() != twoPairRank {
		t.Error()
	}
	if !cardsMakePair(twoPair.highPair, Five, 2) || !cardsMakePair(twoPair.lowPair, Four, 2) {
		t.Error()
	}
	if twoPair.kicker.value != Ace {
		t.Error()
	}
}

func TestPair(t *testing.T) {
	cards := buildCards([]string{"2C", "3H", "4S", "6H", "AH", "AS", "KC"})
	hand, _ := SolveHand(cards)
	pair := hand.(Pair)

	expectedKickers := buildCards([]string{"KC", "5H", "4S"})
	if pair.rank() != pairRank {
		t.Error()
	}
	if !cardsMakePair(pair.pair, Ace, 2) {
		t.Error()
	}
	if !reflect.DeepEqual(expectedKickers, pair.pair) {
		t.Error()
	}
}

func TestInvalidHandSize(t *testing.T) {
	cards := buildCards([]string{"4H", "7S", "8S", "9D"})
	_, err := SolveHand(cards)
	if err == nil {
		t.Error()
	}

	cards = buildCards([]string{"AS", "AC", "AD", "AH", "TD", "8D", "9D", "KD"})
	_, err = SolveHand(cards)
	if err == nil {
		t.Error()
	}
}

func TestCompareDiffHands(t *testing.T) {
	pair, _ := SolveHand(buildCards([]string{"4H", "5C", "8S", "2D", "KS", "8C", "AD"}))
	straight, _ := SolveHand(buildCards([]string{"4H", "5C", "8S", "2D", "KS", "3C", "6D"}))

	if CompareHands(pair, straight) <= 0 {
		t.Error()
	}
	if CompareHands(straight, pair) >= 0 {
		t.Error()
	}
}

func TestCompareStraightFlush(t *testing.T) {
	better, _ := SolveHand(buildCards([]string{"9H", "TH", "JH", "QH", "KH", "AS", "AH"}))
	worse, _ := SolveHand(buildCards([]string{"9H", "TH", "JH", "QH", "KH", "2S", "8H"}))

	if CompareHands(worse, better) <= 0 {
		t.Error()
	}
}

func TestCompareAceWraparound(t *testing.T) {
	better, _ := SolveHand(buildCards([]string{"2H", "3H", "4H", "5H", "9C", "6H", "7C"}))
	worse, _ := SolveHand(buildCards([]string{"2H", "3H", "4H", "5H", "9C", "AH", "8S"}))

	if CompareHands(worse, better) <= 0 {
		t.Error()
	}
}

func TestCompareFourKind(t *testing.T) {
	h1, _ := SolveHand(buildCards([]string{"4H", "4S", "4D", "4C", "9D", "TS"}))
	h2, _ := SolveHand(buildCards([]string{"4H", "4S", "4D", "4C", "3D", "TH"}))

	if CompareHands(h1, h2) != 0 {
		t.Error()
	}

	h1, _ = SolveHand(buildCards([]string{"9H", "9S", "9D", "9C", "3D", "TH"}))

	if CompareHands(h1, h2) <= 0 {
		t.Error()
	}
}

func TestCompareFullHouse(t *testing.T) {
	h1, _ := SolveHand(buildCards([]string{"4H", "4S", "4D", "9C", "9D"}))
	h2, _ := SolveHand(buildCards([]string{"4H", "4S", "4D", "TC", "TD"}))
	h3, _ := SolveHand(buildCards([]string{"5H", "5S", "5D", "9C", "9D"}))

	if CompareHands(h1, h2) <= 0 || CompareHands(h1, h3) <= 0 {
		t.Error()
	}
}

func TestCompareFlush(t *testing.T) {
	h1, _ := SolveHand(buildCards([]string{"4H", "5H", "6H", "7H", "9H"}))
	h2, _ := SolveHand(buildCards([]string{"4H", "5H", "6H", "7H", "AH"}))

	if CompareHands(h1, h2) <= 0 {
		t.Error()
	}
}

func TestCompareStraight(t *testing.T) {
	h1, _ := SolveHand(buildCards([]string{"4C", "5H", "6S", "7D", "8H"}))
	h2, _ := SolveHand(buildCards([]string{"5H", "6S", "7D", "8H", "9S"}))

	if CompareHands(h1, h2) <= 0 {
		t.Error()
	}
}

func TestCompareThreeKind(t *testing.T) {
	h1, _ := SolveHand(buildCards([]string{"4C", "4H", "4S", "7D", "9S"}))
	h2, _ := SolveHand(buildCards([]string{"5H", "5S", "5D", "8H", "9S"}))
	h3, _ := SolveHand(buildCards([]string{"4H", "4S", "4D", "8H", "9S"}))

	if CompareHands(h1, h2) <= 0 || CompareHands(h1, h3) <= 0 {
		t.Error()
	}
}

func TestCompareTwoPair(t *testing.T) {
	h1, _ := SolveHand(buildCards([]string{"4C", "4H", "6S", "6D", "9S"}))
	h2, _ := SolveHand(buildCards([]string{"4H", "4S", "6D", "6H", "TS"}))
	h3, _ := SolveHand(buildCards([]string{"4H", "4S", "7D", "7H", "9S"}))
	h4, _ := SolveHand(buildCards([]string{"5H", "5S", "6D", "6H", "9S"}))

	if CompareHands(h1, h2) <= 0 || CompareHands(h1, h3) <= 0 || CompareHands(h1, h4) <= 0 {
		t.Error()
	}
}

func TestComparePair(t *testing.T) {
	h1, _ := SolveHand(buildCards([]string{"4C", "4H", "6S", "7D", "9S"}))
	h2, _ := SolveHand(buildCards([]string{"4H", "4S", "6D", "8H", "9S"}))
	h3, _ := SolveHand(buildCards([]string{"5H", "5S", "7D", "8H", "9S"}))

	if CompareHands(h1, h2) <= 0 || CompareHands(h1, h3) <= 0 {
		t.Error()
	}
}

func TestCompareHighCard(t *testing.T) {
	h1, _ := SolveHand(buildCards([]string{"2C", "4H", "6S", "7D", "9S"}))
	h2, _ := SolveHand(buildCards([]string{"3C", "4S", "6D", "7H", "9S"}))
	h1Chop, _ := SolveHand(buildCards([]string{"2C", "4S", "6C", "7D", "9H"}))

	if CompareHands(h1, h2) <= 0 || CompareHands(h1, h1Chop) != 0 {
		t.Error()
	}
}
