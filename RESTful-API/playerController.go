package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"

	"net/http"

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
		switch dberr {
		case errors.New("Create Failed fam"):
			respondWithError(w, http.StatusBadRequest, dberr.Error())
		case errors.New("Query Error"):
			respondWithError(w, http.StatusBadRequest, dberr.Error())
		case errors.New("Abnormal number of creates"):
			respondWithError(w, http.StatusNotImplemented, dberr.Error())
		default:
			respondWithError(w, http.StatusBadRequest, dberr.Error())
		}
		return
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
	if err := QuerySearchPlayer(a.db, &p); err != nil {
		switch err {
		case errors.New("Empty"):
			respondWithError(w, http.StatusBadRequest, "Empty")
		case errors.New("Query Error"):
			respondWithError(w, http.StatusBadRequest, "Bad Query")
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Player not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
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
		respondWithError(w, http.StatusInternalServerError, err.Error())
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
		respondWithError(w, http.StatusNotModified, dberr.Error())
		println(dberr.Error())
		return
	}

	//Returning modified object
	respondWithJSON(w, http.StatusOK, p)

}

/*********Game Routes*********/

func (a *App) registerPlayer(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (a *App) loginPlayer(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//accessToken := vars["access_token"]

}

func (a *App) statsPlayer(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (a *App) checkToken(w http.ResponseWriter, r *http.Request) {
	//TODO
}

/*********Helpers*********/

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
