package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createTestRequest(method string, val interface{}) *http.Request {
	json, err := json.Marshal(val)
	if err != nil {
		panic("Failed to marshall value")
	}

	// Since we hook this right up to the handler, there is no need for the uri
	req := httptest.NewRequest(method, "/stub-uri", bytes.NewReader(json))
	return req
}

func readResponse(r *http.Response, dst interface{}) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, dst); err != nil {
		panic(err)
	}

	return
}

func TestAddPlayerApi(t *testing.T) {
	gameReq := createGameRequest{Passphrase: "foobar", StarterChips: 100, BlindSize: 10}
	req := createTestRequest("POST", gameReq)

	recorder := httptest.NewRecorder()
	gameManager := NewGameManager()
	handler := http.HandlerFunc(gameManager.CreateGame)
	handler.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Error()
	}

	gameResponse := getGameResponse{}
	readResponse(recorder.Result(), &gameResponse)

	playerReq := addPlayerRequest{Name: "foo", Seat: 3, Passphrase: gameReq.Passphrase, GameID: gameResponse.GameID}
	req = createTestRequest("POST", playerReq)
	recorder = httptest.NewRecorder()
	handler = http.HandlerFunc(gameManager.AddPlayer)
	handler.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Error()
	}

	playerResponse := addPlayerResponse{}
	readResponse(recorder.Result(), &playerResponse)

	if playerResponse.Player.Chips != gameReq.StarterChips {
		t.Error()
	}
	if playerResponse.Player.Seat != playerReq.Seat {
		t.Error()
	}
}
