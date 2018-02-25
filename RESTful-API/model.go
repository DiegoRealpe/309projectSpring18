package main

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

//Player ...
type Player struct {
	ID          string `json:"id,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	GamesPlayed string `json:"gamesplayed,omitempty"`
	GoalsScored string `json:"goalsscored,omitempty"`
}

//QueryDeleteUser Clears user in database
func QueryDeleteUser(db *sql.DB, IDConstrain string) error {
	return errors.New("Not ready")
	/*var request string
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
	}*/
}

//QueryCreateUser inserts new user in database
func QueryCreateUser(db *sql.DB, NewID string) error {
	return errors.New("Not ready")
	/*var request string
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
	}*/

}

//QueryAllUsers Returns all the users stored in the Clients table
func QueryAllUsers(db *sql.DB) error {
	return errors.New("Not ready")
	/*var rows, err = db.Query("SELECT * FROM Clients")
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
	fmt.Println("Users in total:", results)*/
}

//QuerySearchUser Looks for user in database
func QuerySearchUser(db *sql.DB, Column string, Constrain string) error {
	return errors.New("Not ready")
	/*var request string
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
	}*/
}

/*
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
}*/
