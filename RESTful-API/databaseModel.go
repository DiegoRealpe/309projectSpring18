package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

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
func QueryDeletePlayer(db *sql.DB, p *Player) error {
	request := fmt.Sprintf(`DELETE FROM Players WHERE ID = '%s'`, p.ID)
	result, err := db.Exec(request)
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
func QueryCreatePlayer(db *sql.DB, p *Player) error {
	request := fmt.Sprintf(`INSERT INTO Players (Nickname, GamesPlayed, GamesWon, GoalsScored, Active)
	VALUES ('%s', '0', '0', '0', '0')`, p.Nickname)
	result, err := db.Exec(request)
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
func QueryAllPlayers(db *sql.DB) error {
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
func QuerySearchPlayer(db *sql.DB, p *Player) error {
	if p.ID == "" {
		return errors.New("Empty")
	}
	request := fmt.Sprintf("SELECT * FROM Players WHERE ID = '%s'", p.ID)
	rows, err := db.Query(request)
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

//QueryUpdatePlayer Searchs for a matching ID and updates based on player values given
//Returns errors in case of not finding the correct ID or getting a wrong value
//MODIFIES Player object to overrw
func QueryUpdatePlayer(db *sql.DB, p *Player) error {
	if p.ID == "" {
		return errors.New("Empty ID")
	}
	var mods []string //Declaring slie of values to change

	//any value of the struct that is non nil is updated
	if p.Nickname != "" {
		mods = append(mods, "Nickname", p.Nickname)
	}

	if p.GamesPlayed != "" {
		mods = append(mods, "GamesPlayed", p.GamesPlayed)
	}

	if p.GoalsScored != "" {
		mods = append(mods, "GoalsScored", p.GoalsScored)
	} //Easily can add more

	mods = append(mods, p.ID) //Appending ID which is gonna be modified

	effect, execErr := db.Exec(prepUpdate(mods))
	if execErr != nil {
		return execErr
	}
	i, _ := effect.RowsAffected()
	if int(i) == 0 {
		return errors.New("Not Modified")
	}

	return nil
}

//QueryCreateFBData inserts player information obtained from graph API
func QueryCreateFBData(db *sql.DB, u *AppUser) error {
	request := fmt.Sprintf(`INSERT INTO FacebookData (FacebookID, PlayerID, FullName, Email)
	VALUES ('%s', '%s', '%s', '%s')`, u.FacebookID, u.ID, u.FullName, u.Email)
	result, err := db.Exec(request)
	if err != nil {
		return errors.New("Query Error")
	}
	affected, err2 := result.RowsAffected()
	if err2 != nil {
		return errors.New("Create Failed fam")
	}
	if affected != int64(1) {
		return errors.New("Abnormal number of creates")
	}
	return nil
}

//QuerySetToken creates a new entry on the applicationToken table
//giving the set ID a corresponding appToken
func QuerySetToken(db *sql.DB, ID string, appToken string) error {
	request := fmt.Sprintf(`INSERT INTO TokenTable (applicationToken, playerID)
	VALUES ('%s', '%s')`, appToken, ID)
	result, err := db.Exec(request)
	if err != nil {
		return errors.New("Query Error")
	}
	affected, err2 := result.RowsAffected()
	if err2 != nil {
		return errors.New("Create Failed fam")
	}
	if affected != int64(1) {
		return errors.New("Abnormal number of creates")
	}
	return nil
}

/*********Helpers*********/

//Helper function that uses a slice of string parameters to prepare an update function for a user
//Parameters are on the form of col1, val1, col2, val2 ... ID of user
//Any less than 3 strings will return an empty string instead
func prepUpdate(parameters []string) string {
	i := len(parameters)
	var j int
	if (i%2) != 1 || i < 2 {
		return "Not enough parameters: " + strconv.Itoa(i)
	}
	stmt := "UPDATE Players SET"
	for j < i-3 {
		stmt += fmt.Sprintf("`%s` = '%s',", parameters[j], parameters[j+1])
		j += 2
	}

	return stmt + fmt.Sprintf("`%s` = '%s' WHERE ID = '%s'", parameters[j], parameters[j+1], parameters[j+2])
}
