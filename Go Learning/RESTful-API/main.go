package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//GetAllPlayers ...
func GetAllPlayers(w http.ResponseWriter, r *http.Request) {

}

//GetPlayer ...
func GetPlayer(w http.ResponseWriter, r *http.Request) {

}

//CreatePlayer ...
func CreatePlayer(w http.ResponseWriter, r *http.Request) {

}

//DeletePlayer ...
func DeletePlayer(w http.ResponseWriter, r *http.Request) {

}

// our main function
func main() {
	fmt.Println("Working")
	router := mux.NewRouter()
	router.HandleFunc("/player", GetAllPlayers).Methods("GET")
	router.HandleFunc("/player/{id}", GetPlayer).Methods("GET")
	router.HandleFunc("/player/{id}", CreatePlayer).Methods("POST")
	router.HandleFunc("/player/{id}", DeletePlayer).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

//Player ...
type Player struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

//Address of the player
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}
