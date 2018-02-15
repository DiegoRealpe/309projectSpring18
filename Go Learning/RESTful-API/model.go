package main

import (
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//Player ...
type Player struct {
	ID          string `json:"id,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	GamesPlayed string `json:"gamesplayed,omitempty"`
	GoalsScored string `json:"goalsscored,omitempty"`
}

//GetAllPlayers ...
func GetAllPlayers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(AllPlayers)
}

//GetPlayer ...
func GetPlayer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range AllPlayers {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Player{})
}

//CreatePlayer ...
func CreatePlayer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user Player
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.ID = params["id"]
	AllPlayers = append(AllPlayers, user)
	json.NewEncoder(w).Encode(AllPlayers)
}

//DeletePlayer ...
func DeletePlayer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, user := range AllPlayers {
		if user.ID == params["id"] {
			AllPlayers = append(AllPlayers[:i], AllPlayers[i+1:]...)
			break
		}
		json.NewEncoder(w).Encode(AllPlayers)
	}
}
