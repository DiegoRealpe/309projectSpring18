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

/*********Routes*********/

func (a *App) getPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	p := Player{ID: strconv.Itoa(id)}
	if err := QuerySearchPlayer(a.db, &p); err != nil {
		switch err {
		case errors.New("Empty"):
			respondWithError(w, http.StatusBadRequest, "Empty")
		case errors.New("Query Error"):
			respondWithError(w, http.StatusBadRequest, "Bad Query")
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) createPlayer(w http.ResponseWriter, r *http.Request) {

	var p Player
	decoder := json.NewDecoder(r.Body) //Passing credentials through http request body
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	err := QueryCreatePlayer(a.db, &p)
	if err != nil {
		switch err {
		case errors.New("Create Failed fam"):
			respondWithError(w, http.StatusBadRequest, "Query Return Error")
		case errors.New("Query Error"):
			respondWithError(w, http.StatusBadRequest, "Bad Query")
		case errors.New("Abnormal number of creates"):
			respondWithError(w, http.StatusNotImplemented, "Abnormal number of creates")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) deletePlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil || id == 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	p := Player{ID: strconv.Itoa(id)}
	if err := QueryDeletePlayer(a.db, &p); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
