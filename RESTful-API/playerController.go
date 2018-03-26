package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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

/*********Game Routes*********/

func (a *App) registerPlayer(w http.ResponseWriter, r *http.Request) {
	//1 get token
	token := r.Header.Get("FacebookToken")
	fmt.Println("1")
	//2 check token with fb and get AppUser struct
	user := getFBUser(token)
	if user.Valid == false {
		respondWithError(w, http.StatusForbidden, "Token Error")
		return
	}
	fmt.Println("2")
	//3 get nickname and create player
	var p Player
	decoder := json.NewDecoder(r.Body) //Passing credentials through http request body
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusNotAcceptable, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	fmt.Println("3 decoding")
	dberr := QueryCreatePlayer(a.db, &p)
	if dberr != nil {
		handleDBErrors(w, dberr)
		return
	}
	fmt.Println("3")
	//4 create facebook data for player
	user.ID = p.ID //giving appuser's info the id of the player it belongs to
	dberr = QueryCreateFBData(a.db, &user)
	if dberr != nil {
		handleDBErrors(w, dberr)
		return
	}
	fmt.Println("4")
	//5 create application token and update table
	dberr = QuerySetToken(a.db, p.ID, appTokenGen(p.ID))
	if dberr != nil {
		handleDBErrors(w, dberr)
		return
	}
	fmt.Println("5")
	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) loginPlayer(w http.ResponseWriter, r *http.Request) {
	//TODO
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

func handleDBErrors(w http.ResponseWriter, dberr error) {
	switch dberr {

	case errors.New("Update Error"):
		respondWithError(w, http.StatusNotModified, dberr.Error())
	case errors.New("Create Failed fam"):
		respondWithError(w, http.StatusBadRequest, dberr.Error())
	case errors.New("Query Error"):
		respondWithError(w, http.StatusBadRequest, dberr.Error())
	case errors.New("Abnormal number of creates"):
		respondWithError(w, http.StatusNotImplemented, dberr.Error())
	case errors.New("Empty"):
		respondWithError(w, http.StatusBadRequest, "Empty")
	case sql.ErrNoRows:
		respondWithError(w, http.StatusNotFound, "Player not found")
	default:
		respondWithError(w, http.StatusInternalServerError, dberr.Error())
	}
	return
}

//Stupid funtion to try and make a random number
func appTokenGen(ID string) string {
	i, _ := strconv.Atoi(ID)
	i = i*186282 + i*299792 //speed of light in mps and kms
	return strconv.Itoa(i)
}
