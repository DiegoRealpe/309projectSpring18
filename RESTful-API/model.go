package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//Player ...
type Player struct {
	ID          string `json:"id,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	GamesPlayed string `json:"gamesplayed,omitempty"`
	GoalsScored string `json:"goalsscored,omitempty"`
}

//QueryDeletePlayer Clears Player in database
func (p *Player) QueryDeletePlayer(db *sql.DB) error {
	request := fmt.Sprintf(`DELETE FROM Players WHERE ID = '%s'`, p.ID)
	var result, err = db.Exec(request)
	if err != nil {
		return errors.New("Query Error")
	}
	affected, err2 := result.RowsAffected()
	if err2 != nil {
		return errors.New("Resulting Rows Error")
	}
	if affected != int64(1) {
		return errors.New("None Found")
	}
	return nil
}

//QueryCreatePlayer inserts new Player in database
func (p *Player) QueryCreatePlayer(db *sql.DB) error {
	request := fmt.Sprintf(`INSERT INTO Players (Nickname, GamesPlayed, GamesWon, GoalsScored, Active)
	VALUES ('%s', '0', '0', '0', '0')`, p.Nickname)
	var result, err = db.Exec(request)
	if err != nil {
		return errors.New("Query Error")
	}
	affected, err2 := result.RowsAffected()
	if err2 != nil {
		return errors.New("Create Failed fam")
	}
	if affected == int64(1) {
		db.QueryRow("SELECT ID FROM Players WHERE Nickname = ?", p.Nickname).Scan(&p.ID)
		p.GamesPlayed = "0"
		p.GoalsScored = "0"
		return nil
	}
	return errors.New("Abnormal number of creates")
}

//QueryAllPlayers Returns all the Players stored in the Players table
func (p *Player) QueryAllPlayers(db *sql.DB) error {
	return errors.New("Not ready")
	/*var rows, err = db.Query("SELECT * FROM Players")
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
	fmt.Println("Players in total:", results)*/
}

//QuerySearchPlayer Looks for Player in database
func (p *Player) QuerySearchPlayer(db *sql.DB) error {
	if p.ID == "" {
		return errors.New("Empty")
	}
	request := fmt.Sprintf("SELECT * FROM Players WHERE ID = '%s'", p.ID)
	var rows, err = db.Query(request)
	if err != nil {
		return errors.New("Query Error")
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
		p.Nickname = Nickname
		p.GamesPlayed = string(a)
		p.GoalsScored = string(c)
	}
	if results == 0 { //Diego from the future, you idiot, dont move this from here
		return sql.ErrNoRows
	}

	return nil
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
	var Player Player
	_ = json.NewDecoder(r.Body).Decode(&Player)
	Player.ID = params["id"]
	AllPlayers = append(AllPlayers, Player)
	json.NewEncoder(w).Encode(AllPlayers)
}

//DeletePlayer ...
func DeletePlayer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, Player := range AllPlayers {
		if Player.ID == params["id"] {
			AllPlayers = append(AllPlayers[:i], AllPlayers[i+1:]...)
			break
		}
		json.NewEncoder(w).Encode(AllPlayers)
	}
}*/
