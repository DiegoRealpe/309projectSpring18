package main

import (
	"encoding/json"
	"errors"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

/*********CRUD Routes*********/

func (a *App) getPlayer(w http.ResponseWriter, r *http.Request) {

	//Veriy if the game service is accessing
	if securityErr := verifyAccess(r); securityErr != nil {
		respondWithError(w, http.StatusUnauthorized, securityErr.Error())
		return
	}
	//Obtaining one value, ID from mux parameters to create player
	vars := mux.Vars(r)
	token := vars["AppToken"]
	id, assertErr := QueryAssertToken(a.db, token)
	if assertErr != nil {
		handleDBErrors(w, assertErr)
		return
	}
	p := Player{ID: id}

	//Executing search query
	err := QuerySearchPlayer(a.db, &p)
	if err != nil {
		handleDBErrors(w, err)
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) deletePlayer(w http.ResponseWriter, r *http.Request) {

	//Veriy if the game service is accessing
	if securityErr := verifyAccess(r); securityErr != nil {
		respondWithError(w, http.StatusUnauthorized, securityErr.Error())
		return
	}

	//Obtaining ID from mux variables
	vars := mux.Vars(r)
	token := vars["AppToken"]
	id, assertErr := QueryAssertToken(a.db, token)
	if assertErr != nil {
		handleDBErrors(w, errors.New("Invalid User Token"))
		return
	}
	p := Player{ID: id}

	//Executing delete query model
	if err := QueryDeletePlayer(a.db, &p); err != nil {
		handleDBErrors(w, err)
		return
	}

	//Updating table to reflect rank
	rankErr := QueryRankTrigger(a.db)
	if rankErr != nil {
		handleDBErrors(w, rankErr)
		return
	}

	respondWithJSON(w, http.StatusAccepted, nil)
}

func (a *App) updatePlayer(w http.ResponseWriter, r *http.Request) {

	//Veriy if the gameservice is accessing
	if securityErr := verifyAccess(r); securityErr != nil {
		respondWithError(w, http.StatusUnauthorized, securityErr.Error())
		return
	}
	//Getting ID from mux parameter
	vars := mux.Vars(r)
	token := vars["AppToken"]
	id, assertErr := QueryAssertToken(a.db, token)
	if assertErr != nil {
		handleDBErrors(w, errors.New("Invalid User Token"))
		return
	}

	//Getting variables to change from http.request.body
	var p Player
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	p.ID = id

	//Executing Query model
	dberr := QueryUpdatePlayer(a.db, &p)
	if dberr != nil {
		handleDBErrors(w, errors.New("Update Error"))
		return
	}

	//Updating table to reflect rank
	rankErr := QueryRankTrigger(a.db)
	if rankErr != nil {
		handleDBErrors(w, rankErr)
		return
	}

	//Returning modified object
	respondWithJSON(w, http.StatusOK, p)

}
