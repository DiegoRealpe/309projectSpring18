package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

/*********Game Routes*********/

func (a *App) registerPlayer(w http.ResponseWriter, r *http.Request) {
	//1 get token
	token := r.Header.Get("FacebookToken")
	//2 check token with fb and get AppUser struct
	user := getFBUser(token)
	if user.Valid == false {
		respondWithError(w, http.StatusForbidden, "Token Error")
		return
	}
	//3 get nickname and create player
	var p Player
	decoder := json.NewDecoder(r.Body) //Passing credentials through http request body
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusNotAcceptable, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	dberr := QueryCreatePlayer(a.db, &p)
	if dberr != nil {
		handleDBErrors(w, dberr)
		return
	}
	//4 create facebook data for player
	user.ID = p.ID //giving appuser's info the id of the player it belongs to
	dberr = QueryCreateFBData(a.db, &user)
	if dberr != nil {
		handleDBErrors(w, dberr)
		return
	}
	//5 create application token and update table
	apptoken := appTokenGen(p.ID)
	dberr = QuerySetToken(a.db, p.ID, apptoken, 1)
	if dberr != nil {
		handleDBErrors(w, dberr)
		return
	}

	profile := PlayerProfile{Profile: p, AppToken: apptoken}
	respondWithJSON(w, http.StatusCreated, profile)
}

func (a *App) loginPlayer(w http.ResponseWriter, r *http.Request) {
	//1 get fb token
	token := r.Header.Get("FacebookToken")
	//2 check the graph api and get AppUser object
	user := getFBUser(token)
	if user.Valid == false {
		respondWithError(w, http.StatusForbidden, "Token Error")
		return
	}
	//3 query AppUser FBID in FBdatatable to get game id
	dberr := QueryGetFBDataID(a.db, &user)
	if dberr != nil {
		handleDBErrors(w, dberr)
		return
	}
	//4 use GET player model to get info
	p := Player{ID: user.ID}
	dberr = QuerySearchPlayer(a.db, &p)
	if dberr != nil {
		handleDBErrors(w, dberr)
		return
	}
	//5 use GET token model to get apptoken (handle expired)
	apptoken, dberr := QueryGetToken(a.db, user.ID)
	if dberr != nil {
		handleDBErrors(w, dberr)
		return
	}
	//6 answer with player struct and apptoken
	profile := PlayerProfile{Profile: p, AppToken: apptoken}
	respondWithJSON(w, http.StatusOK, profile)
}

func (a *App) statsPlayer(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (a *App) checkToken(w http.ResponseWriter, r *http.Request) {
	//TODO
}

/*********Helpers*********/

//Stupid funtion to try and make a random number
func appTokenGen(ID string) string {
	i, _ := strconv.Atoi(ID)
	i = i*186282 + i*299792 //speed of light in mps and kms
	return strconv.Itoa(i)
}
