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
			panic("Unkown card calue")
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

func TestHighCardBasic(t *testing.T) {
	cards := buildCards([]string{"2H", "4C", "8D", "TS", "AS", "3H", "9C"})
	highCard := SolveHand(cards).(*HighCard)

	expected := Card{suit: Spades, value: Ace}
	if !reflect.DeepEqual(highCard.highCard, expected) {
		t.Errorf("%+v\n != %+v\n", highCard.highCard, expected)
	}

	if highCard.rank() != highCardRank {
		t.Error()
	}
}

func TestFlushBasic(t *testing.T) {
	cards := buildCards([]string{"AS", "2S", "4S", "8S", "TS"})
	flush := SolveHand(cards).(*Flush)

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
	flush := SolveHand(cards).(*Flush)

	expected := buildCards([]string{"AS", "TS", "9S", "8S", "4S"})
	if !reflect.DeepEqual(expected, flush.sortedCards) {
		t.Error()
	}
}

func TestStraightBasic(t *testing.T) {
	/*
		cards := buildCards([]string{"2H", "3C", "6S", "TS", "QC", "4D", "5D"})
		straight := SolveHand(cards).(*Straight)

		expected := buildCards([]string{"6S", "5D", "4D", "3C", "2H"})
		if !reflect.DeepEqual(expected, straight.sortedCards) {
			t.Error()
		}
	*/
}

func TestStraightDuplicateCards(t *testing.T) {

}

func TestStraightLong(t *testing.T) {

}

func TestStraightAceLow(t *testing.T) {

}

func TestStraightAceHigh(t *testing.T) {

}

func TestStraightAceWraparound(t *testing.T) {

}
