package api

import (
	"net/http"
	"sync"

	"github.com/brian-a-esch/httpoker/poker"
)

// GameManager is the top level object which exposes some http endpoints
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

	game, ok := manager.games[get.GameID]
	if !ok {
		http.Error(w, "Could not find/load game", http.StatusBadRequest)
		return
	}

	passprhase, ok := manager.gamePassphrases[get.GameID]
	if !ok {
		http.Error(w, "Could not find/load game", http.StatusBadRequest)
		return
	}

	if get.Passphrase != passprhase {
		http.Error(w, "Could not find/load game", http.StatusBadRequest)
		return
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
