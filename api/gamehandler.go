package api

import (
	"errors"
	"net/http"
	"sync"

	"github.com/brian-a-esch/httpoker/poker"
)

// GameManager is the top level object which exposes some http endpoints. There is a mutex
// since there is shared state which needs to be managed. Reader writer locks were considered,
// and if this were to be used for more than a few games, then it would probably be implemented.
// But since typically we only play one game at a time, having one lock for all the games is fine
// and it avoids having to lock at the game level
type GameManager struct {
	sync.Mutex
	gameCounter     int
	games           map[int]poker.Game
	gamePassphrases map[int]string
}

// NewGameManager allocates a new GameManger
func NewGameManager() GameManager {
	return GameManager{gameCounter: 0, games: make(map[int]poker.Game), gamePassphrases: make(map[int]string)}
}

type createGameRequest struct {
	Passphrase   string `json:"passphrase"`
	StarterChips int    `json:"starterChips"`
	BlindSize    int    `json:"blindSize"`
}

type getGameRequest struct {
	GameID     int    `json:"gameID"`
	Passphrase string `json:"passphrase"`
}

// TODO should we just have game provide a "serialize" method?
type getGameResponse struct {
	GameID       int                  `json:"gameID"`
	StarterChips int                  `json:"starterChips"`
	BlindSize    int                  `json:"blindSize"`
	EmptySeats   []int                `json:"emptySeats"`
	Players      map[int]poker.Player `json:"players"`
}

// Game is a restful endpoint for getting a poker game
func (manager *GameManager) Game(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Must be POST", http.StatusMethodNotAllowed)
	}

	get := getGameRequest{}
	if ok := decodeJSONBody(w, r, &get); !ok {
		return
	}

	manager.Lock()
	defer manager.Unlock()

	game, err := manager.resolveGame(get.GameID, get.Passphrase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	response := getGameResponse{
		GameID:       get.GameID,
		StarterChips: game.StarterChips(),
		BlindSize:    game.BlindSize(),
		EmptySeats:   game.EmptySeats(),
		Players:      game.Players(),
	}
	sendJSONResponse(w, response)
}

// CreateGame creates a game in the GameManager
func (manager *GameManager) CreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Must be POST", http.StatusMethodNotAllowed)
	}

	create := createGameRequest{}
	if ok := decodeJSONBody(w, r, &create); !ok {
		return
	}

	game, err := poker.NewGame(create.StarterChips, create.BlindSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	manager.Lock()
	defer manager.Unlock()

	gameID := manager.gameCounter
	manager.gameCounter++
	if _, ok := manager.games[gameID]; ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	manager.games[gameID] = game
	manager.gamePassphrases[gameID] = create.Passphrase

	response := getGameResponse{GameID: gameID, StarterChips: game.StarterChips(), BlindSize: game.BlindSize()}
	sendJSONResponse(w, response)
}

type addPlayerRequest struct {
	Name       string `json:"name"`
	Passphrase string `json:"passphrase"`
	Seat       int    `json:"seat"`
	GameID     int    `json:"gameID"`
}

type addPlayerResponse struct {
	Player poker.Player `json:"player"`
	Secret int          `json:"secret"`
}

// AddPlayer adds a player to the game
func (manager *GameManager) AddPlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Must be POST", http.StatusMethodNotAllowed)
	}

	addRequest := addPlayerRequest{}
	if ok := decodeJSONBody(w, r, &addRequest); !ok {
		return
	}

	manager.Lock()
	defer manager.Unlock()

	game, err := manager.resolveGame(addRequest.GameID, addRequest.Passphrase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	player, err := game.AddPlayer(addRequest.Name, addRequest.Seat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	sendJSONResponse(w, addPlayerResponse{
		Player: player,
		Secret: player.Secret(),
	})

}

func (manager *GameManager) resolveGame(gameID int, passphrase string) (*poker.Game, error) {
	// The error functions here are intentionally obtuse
	game, ok := manager.games[gameID]
	if !ok {
		return nil, errors.New("Could not find/load game")
	}

	foundPassphrase, ok := manager.gamePassphrases[gameID]
	if !ok {
		return nil, errors.New("Could not find/load game")
	}

	if passphrase != foundPassphrase {
		return nil, errors.New("Could not find/load game")
	}

	return &game, nil
}
