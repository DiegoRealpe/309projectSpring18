package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

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

//QueryAllUsers Returns all the users stored in the Clients table
func QueryAllUsers(db *sql.DB) {
	var rows, err = db.Query("SELECT * FROM Clients")
	if err != nil {
		fmt.Println(err)
	}

	var ID string
	var Nickname string
	var a, b, c, d int
	for rows.Next() {
		err := rows.Scan(&ID, &Nickname, &a, &b, &c, &d)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("ID = ", ID, "Nickname = ", Nickname)
	}
}

//QuerySearchUser Looks for user in database
func QuerySearchUser(db *sql.DB, IDConstrain string) {
	var request string
	request = fmt.Sprintf("SELECT * FROM Clients WHERE ID = '%s'", IDConstrain)
	var rows, err = db.Query(request)
	if err != nil {
		fmt.Println(err)
	}

	var ID, Nickname string
	var results, a, b, c, d int

	for rows.Next() {
		err := rows.Scan(&ID, &Nickname, &a, &b, &c, &d)
		if err != nil {
			fmt.Println(err)
		}
		results++
		fmt.Println("ID = ", ID, "Nickname = ", Nickname)
	}
	if results == 0 {
		fmt.Println("No Results Found!")
	}
}

// our main function
func main() {
	fmt.Println("Working")

	db, err := sql.Open("mysql",
		"dbu309mg6:1XFA40wc@tcp(mysql.cs.iastate.edu:3306)/db309mg6")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	QuerySearchUser(db, "0")

	/*
		db.Query(`INSERT INTO Clients(ID, Nickname, GamesPlayed, GamesWon, GoalsScored, Online)
				  VALUES(1, Ryan, 0, 0, 0, 0)`)
	*/

	return
	fmt.Println("Server")
	AllPlayers = append(AllPlayers, Player{ID: "1", Nickname: "Nolan", GamesPlayed: "0", GoalsScored: "0"})
	AllPlayers = append(AllPlayers, Player{ID: "2", Nickname: "Diego", GamesPlayed: "5", GoalsScored: "10"})
	router := mux.NewRouter()
	router.HandleFunc("/player", GetAllPlayers).Methods("GET")
	router.HandleFunc("/player/{id}", GetPlayer).Methods("GET")
	router.HandleFunc("/player/{id}", CreatePlayer).Methods("POST")
	router.HandleFunc("/player/{id}", DeletePlayer).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

//Player ...
type Player struct {
	ID          string `json:"id,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	GamesPlayed string `json:"gamesplayed,omitempty"`
	GoalsScored string `json:"goalsscored,omitempty"`
}

//AllPlayers complete list
var AllPlayers []Player

/*db.Query(`CREATE DATABASE Diego_Test;
USE Diego_Test;
CREATE TABLE Clients (
	ID UNIQUEIDENTIFIER PRIMARY KEY,
	Nickname VARCHAR(50) NOT NULL,
	GamesPlayed INT NOT NULL,
	GamesWon INT NOT NULL,
	GoalsScored INT NOT NULL,
	Online INT NOT NULL)`)*/
