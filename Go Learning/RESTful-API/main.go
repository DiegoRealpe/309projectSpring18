package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//QueryDeleteUser Clears user in database
func QueryDeleteUser(db *sql.DB, IDConstrain string) {
	var request string
	request = fmt.Sprintf(`DELETE FROM Clients WHERE ID = '%s'`, IDConstrain)
	var result, err = db.Exec(request)
	if err != nil {
		fmt.Println(err)
	}
	affected, err2 := result.RowsAffected()
	if err2 != nil {
		fmt.Println(err2)
	}
	if affected == int64(1) {
		fmt.Println("Deleted User:", IDConstrain)
	} else {
		fmt.Println("None Found")
	}
}

//QueryCreateUser inserts new user in database
func QueryCreateUser(db *sql.DB, NewID string) {
	var request string
	request = fmt.Sprintf(`INSERT INTO Clients (Nickname, GamesPlayed, GamesWon, GoalsScored, Active) VALUES ('%s', '0', '0', '0', '0')`, NewID)
	var result, err = db.Exec(request)
	if err != nil {
		fmt.Println(err)
	}
	affected, err2 := result.RowsAffected()
	if err2 != nil {
		fmt.Println(err2)
	}
	if affected == int64(1) {
		fmt.Println("Created New User:", NewID)
	} else {
		fmt.Println("Create Failed")
	}

}

//QueryAllUsers Returns all the users stored in the Clients table
func QueryAllUsers(db *sql.DB) {
	var rows, err = db.Query("SELECT * FROM Clients")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var ID string
	var Nickname string
	var results, a, b, c, d int
	for rows.Next() {
		err := rows.Scan(&ID, &Nickname, &a, &b, &c, &d)
		if err != nil {
			fmt.Println(err)
		}
		results++
		fmt.Println("ID = ", ID, "Nickname = ", Nickname)
	}
	fmt.Println("Users in total:", results)
}

//QuerySearchUser Looks for user in database
func QuerySearchUser(db *sql.DB, Column string, Constrain string) {
	var request string
	request = fmt.Sprintf("SELECT * FROM Clients WHERE %s = '%s'", Column, Constrain)
	var rows, err = db.Query(request)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

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

func main() {
	fmt.Println("Working")

	a := App{}
	a.Initialize()

	a.Run(":8080")

	QuerySearchUser(db, "Nickname", "Nolan")

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

//AllPlayers complete list
var AllPlayers []Player

/*db.Query(`CREATE DATABASE Diego_Test;
USE Diego_Test;
CREATE TABLE Clients (
	ID INT PRIMARY KEY,
	Nickname VARCHAR(50) NOT NULL,
	GamesPlayed INT NOT NULL,
	GamesWon INT NOT NULL,
	GoalsScored INT NOT NULL,
	Online INT NOT NULL)`)*/
