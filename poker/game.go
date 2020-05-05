package poker

import (
	"fmt"
)

// NumSeats is the number of seats in a game
const NumSeats int = 8

// Player is a registered participant in the game
type Player struct {
	Name   string `json:"name"`
	Chips  int    `json:"chips"`
	Seat   int    `json:"seat"`
	secret int
}

// Secret is a getter for the player secret. Should ONLY be used to give access to
// a user who actually has access to this player in the game.
func (player *Player) Secret() int {
	return player.secret
}

// Game is the highest level object, representing and entire poker game
type Game struct {
	players      map[int]Player
	deck         Deck
	starterChips int
	blindSize    int
	button       int
}

// NewGame creates a game with the specified rules
func NewGame(starterChips int, blindSize int) (Game, error) {
	if starterChips <= 0 {
		return Game{}, fmt.Errorf("Game needs to have a positive chip count")
	}
	if blindSize <= 0 || blindSize%2 != 0 {
		return Game{}, fmt.Errorf("Game needs to have a positive blind size, which is divisible by 2")
	}

	return Game{
		players:      make(map[int]Player),
		deck:         NewDeck(),
		starterChips: starterChips,
		blindSize:    blindSize,
		button:       -1,
	}, nil
}

// AddPlayer adds a new player to the game and generates some values for them
// Returns an error for invalid arguments or the game being full
func (game *Game) AddPlayer(name string, seat int) (Player, error) {
	if len(game.players) >= NumSeats {
		return Player{}, fmt.Errorf("Game already at capacity of %d players", NumSeats)
	}

	if seat < 0 || seat > 8 {
		return Player{}, fmt.Errorf("Invalid seat number %d", seat)
	}

	if _, ok := game.players[seat]; ok {
		return Player{}, fmt.Errorf("Already have player at seat %d", seat)
	}

	secret := int(seat) * -1
	player := Player{Name: name, Chips: game.starterChips, Seat: seat, secret: secret}
	game.players[seat] = player
	return player, nil
}

// Players gets the current players in the game
func (game *Game) Players() map[int]Player {
	return game.players
}

// EmptySeats returns a slice of the empty seats
func (game *Game) EmptySeats() []int {
	result := make([]int, 0, NumSeats-len(game.players))
	for i := 0; i < NumSeats; i++ {
		if _, ok := game.players[i]; !ok {
			result = append(result, i)
		}
	}

	return result
}

// StarterChips get the starting chip size
func (game *Game) StarterChips() int {
	return game.starterChips
}

// BlindSize gets the blind size
func (game *Game) BlindSize() int {
	return game.blindSize
}
