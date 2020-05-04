package poker

import (
	"testing"
)

func TestGameInvalidSeat(t *testing.T) {
	game, _ := NewGame(100, 5)
	_, err := game.AddPlayer("Foo", -1)
	if err == nil {
		t.Error()
	}

	_, err = game.AddPlayer("bar", 8)
	if err == nil {
		t.Error()
	}
}
