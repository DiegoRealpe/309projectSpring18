package main

import (
	"database/sql"
	"errors"
	"strconv"
	"time"
)

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
	if AppToken == "" {
		return "", errors.New("Empty AppToken")
	}
	row, err := db.Query(`SELECT ID, expiration FROM TokenTable 
	JOIN Players ON TokenTable.playerID = Players.ID WHERE applicationToken = ?`, AppToken)
	if err != nil {
		return "", err
	}
	var ID, exp int64
	row.Next()
	err = row.Scan(&ID, &exp)
	if err != nil {
		return "", errors.New("Player Not Found")
	}
	if exp < time.Now().Unix() {
		return "", errors.New("Application Token Expired" + err.Error())
	}
	return strconv.Itoa(int(ID)), nil
}
