package main

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//Player ...
type Player struct {
	ID          string `json:"id,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	GamesPlayed string `json:"gamesplayed,omitempty"`
	GamesWon    string `json:"gameswon,omitempty"`
	GoalsScored string `json:"goalsscored,omitempty"`
	RankWin     string `json:"rankwin,omitempty"`
	RankScore   string `json:"rankscore,omitempty"`
}

//PlayerProfile contains a player struct with its respective apptoken
type PlayerProfile struct {
	Profile  Player
	AppToken string `json:"ApplicationToken,omitempty"`
	Error    string `json:"error-message,omitempty"`
}

//QueryDeletePlayer Clears Player in database
func QueryDeletePlayer(db *sql.DB, p *Player) error {
	result, err := db.Exec(`DELETE FROM Players WHERE ID = ?`, p.ID)
	if err != nil {
		return err
	}
	affected, err2 := result.RowsAffected()
	if err2 != nil {
		return errors.New("Rows Affected")
	}
	if affected != int64(1) {
		return errors.New("Player Not Found")
	}
	return nil
}

//QueryCreatePlayer inserts new Player in database
func QueryCreatePlayer(db *sql.DB, p *Player) error {
	result, err := db.Exec(`INSERT INTO Players (Nickname)
	VALUES (?)`, p.Nickname)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == int64(0) {
		return errors.New("Create Fail")
	}
	db.QueryRow("SELECT ID FROM Players WHERE Nickname = ?", p.Nickname).Scan(&p.ID)
	p.GamesPlayed = "0"
	p.GoalsScored = "0"
	return nil
}

//QuerySearchPlayer Looks for Player in database
func QuerySearchPlayer(db *sql.DB, p *Player) error {
	if p.ID == "" {
		return errors.New("Invalid user ID")
	}

	rows, err := db.Query(`SELECT * FROM Players WHERE ID = ?`, p.ID)
	if err != nil {
		return errors.New("Select Failed" + err.Error())
	}
	defer rows.Close()

	var ID, Nickname string
	var results, gamesplayed, gameswon, goalsscored, rankwin, rankscore int

	for rows.Next() {
		err2 := rows.Scan(&ID, &Nickname, &gamesplayed, &gameswon, &goalsscored, &rankwin, &rankscore)
		if err2 != nil {
			return errors.New("Scan Rows Failed" + err2.Error())
		}

		results++
		p.Nickname = Nickname
		p.GamesPlayed = strconv.Itoa(gamesplayed)
		p.GamesWon = strconv.Itoa(gameswon)
		p.GoalsScored = strconv.Itoa(goalsscored)
		p.RankWin = strconv.Itoa(rankwin)
		p.RankScore = strconv.Itoa(rankscore)
	}
	if results == 0 { //Diego from the future, you idiot, dont move this from here
		return sql.ErrNoRows
	}
	return nil
}

//QueryUpdatePlayer Searchs for a matching ID and updates based on player values given
//Returns errors in case of not finding the correct ID or getting a wrong value
//MODIFIES Player object to overrwrite
func QueryUpdatePlayer(db *sql.DB, p *Player) error {
	if p.ID == "" {
		return errors.New("Invalid user ID")
	}
	var mods []string //Declaring slice of values to change

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
	i, err := effect.RowsAffected()
	if err != nil {
		return err
	}
	if int(i) == 0 {
		return errors.New("Not Modified")
	}

	return nil
}

//QueryCreateFBData inserts player information obtained from graph API
func QueryCreateFBData(db *sql.DB, u *AppUser) error {
	result, err := db.Exec(`INSERT INTO FacebookData (FacebookID, PlayerID, FullName, Email)
	VALUES (?, ?, ?, ?)`, u.FacebookID, u.ID, u.FullName, u.Email)
	if err != nil {
		return errors.New("Get FB Data ID Error - Execution" + err.Error())
	}
	affected, err2 := result.RowsAffected()
	if err2 != nil {
		return errors.New("Get FB Data ID Error - Result" + err2.Error())
	}
	if affected == int64(0) {
		return errors.New("Create Fail")
	}
	return nil
}

//QueryGetFBDataID Looks in the db if there is an ID corresponding the AppUser's FB ID
func QueryGetFBDataID(db *sql.DB, u *AppUser) error {
	row, err := db.Query("SELECT PlayerID FROM FacebookData WHERE FacebookID = ?", u.FacebookID)
	if err != nil {
		return err
	}
	row.Next()
	err2 := row.Scan(&u.ID)
	if err2 != nil {
		return errors.New(u.ID)
	}
	return nil
}

//QuerySetToken creates a new entry on the applicationToken table
//giving the set ID a corresponding appToken and assigning an expiration
func QuerySetToken(db *sql.DB, ID string, appToken string, tokenLife int) error {
	result, err := db.Exec(`INSERT INTO TokenTable (applicationToken, playerID, expiration)
	VALUES (?, ?, ?)`, appToken, ID, getExpiration(tokenLife))
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected != int64(1) {
		return errors.New("Create Fail")
	}
	return nil
}

//QueryGetUpdateToken query to return the token assigned to a specific ID
func QueryGetUpdateToken(db *sql.DB, ID string) (string, error) {
	res, err := db.Exec(`UPDATE TokenTable SET expiration = ? WHERE playerID = ?`, getExpiration(2), ID)
	if err != nil {
		return "", err
	}
	affected, _ := res.RowsAffected()
	if int(affected) == 0 {
		return "", sql.ErrNoRows
	}

	row, err := db.Query("SELECT applicationToken FROM TokenTable WHERE playerID = ?",
		ID)
	if err != nil {
		return "", err
	}
	var tok string
	row.Next()
	err = row.Scan(&tok)
	if err != nil {
		return "", errors.New("Player Not Found")
	}
	return tok, nil
}

//QueryGetToken query to return the token assigned to a specific ID
func QueryGetToken(db *sql.DB, ID string) (string, error) {
	row, err := db.Query("SELECT applicationToken, expiration FROM TokenTable WHERE playerID = ?", ID)
	if err != nil {
		return "", err
	}
	var tok string
	var exp int64
	row.Next()
	err = row.Scan(&tok, &exp)
	if err != nil {
		return "", errors.New("Player Not Found")
	}
	if exp < time.Now().Unix() {
		return "", errors.New("Application Token Expired")
	}
	return tok, nil
}

//QueryAssertToken returns the nickname of the given apptoken or 404
func QueryAssertToken(db *sql.DB, AppToken string) (string, error) {
	row, err := db.Query(`SELECT Nickname, expiration FROM TokenTable 
	JOIN Players ON TokenTable.playerID = Players.ID WHERE applicationToken = ?`, AppToken)
	if err != nil {
		return "", err
	}
	var Nickname string
	var exp int64
	row.Next()
	err = row.Scan(&Nickname, &exp)
	if err != nil {
		return "", errors.New("Player Not Found")
	}
	if exp < time.Now().Unix() {
		return "", errors.New("Application Token Expired")
	}
	return Nickname, nil
}
