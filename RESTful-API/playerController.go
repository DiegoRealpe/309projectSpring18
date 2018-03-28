package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

/*********CRUD Routes*********/

func (a *App) createPlayer(w http.ResponseWriter, r *http.Request) {

	//Obtaining specifications through json body
	var p Player
	decoder := json.NewDecoder(r.Body) //Passing credentials through http request body
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusNotAcceptable, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	//Executing Create model
	dberr := QueryCreatePlayer(a.db, &p)
	if dberr != nil {
		handleDBErrors(w, dberr)
	}
	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) getPlayer(w http.ResponseWriter, r *http.Request) {
	//Obtaining one value, ID from mux parameters to create player
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	p := Player{ID: strconv.Itoa(id)}

	//Executing search query
	err = QuerySearchPlayer(a.db, &p)
	if err != nil {
		handleDBErrors(w, err)
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) deletePlayer(w http.ResponseWriter, r *http.Request) {

	//Obtaining ID from mux variables
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil || id == 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	p := Player{ID: strconv.Itoa(id)}

	//Executing delete query model
	if err := QueryDeletePlayer(a.db, &p); err != nil {
		handleDBErrors(w, errors.New("Delete Player Query Error"))
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) updatePlayer(w http.ResponseWriter, r *http.Request) {

	//Getting ID from mux parameter
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil || id == 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	//Getting variables to change from http.request.body
	var p Player
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	p.ID = strconv.Itoa(id)

	//Executing Query model
	dberr := QueryUpdatePlayer(a.db, &p)
	if dberr != nil {
		handleDBErrors(w, errors.New("Update Error"))
		return
	}

	//Returning modified object
	respondWithJSON(w, http.StatusOK, p)

}
