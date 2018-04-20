package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

/*********Game Routes*********/

func (a *App) registerPlayer(w http.ResponseWriter, r *http.Request) {
	//1 get token
	token := r.Header.Get("FacebookToken")
	//2 check token with fb and get AppUser struct
	user := getFBUser(token)
	if user.Valid == false {
		respondWithError(w, http.StatusConflict, "Facebook Token Error")
		return
	}
	//3 get nickname and create player
	var p Player
	decoder := json.NewDecoder(r.Body) //Passing credentials through http request body
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	nickLength := len(p.Nickname)
	if nickLength > 15 || nickLength < 1 {
		respondWithError(w, http.StatusNotImplemented, "Nickname Length Error")
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

	//Updating table to reflect rank
	rankErr := QueryRankTrigger(a.db)
	if rankErr != nil {
		handleDBErrors(w, rankErr)
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
		respondWithError(w, http.StatusConflict, "Facebook Token Error")
		return
	}
	//3 query AppUser FBID in FBdatatable to get game id
	dberr := QueryGetFBDataID(a.db, &user)
	if dberr != nil {
		respondWithError(w, http.StatusNotFound, "Associated FaceBook ID "+
			dberr.Error()+" not registered in DB records")
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
	apptoken, dberr := QueryGetUpdateToken(a.db, user.ID)
	if dberr != nil {
		fmt.Println(dberr.Error())
		handleDBErrors(w, errors.New("Get Apptoken Error"))
		return
	}
	//6 answer with player struct and apptoken
	profile := PlayerProfile{Profile: p, AppToken: apptoken}
	respondWithJSON(w, http.StatusAccepted, profile)
}

func (a *App) statsPlayer(w http.ResponseWriter, r *http.Request) {
	//1 obtain ID
	FacebookID := r.Header.Get("FacebookID")
	//2 Look for the stats in the database
	profile := QueryProfileFromToken(a.db, FacebookID)
	if profile.Error != "" {
		respondWithError(w, http.StatusInternalServerError, profile.Error)
	}
	//3 send struct of stats
	respondWithJSON(w, http.StatusOK, profile)
}

func (a *App) checkToken(w http.ResponseWriter, r *http.Request) {

	//Veriy if the game service is accessing
	if securityErr := verifyAccess(r); securityErr != nil {
		respondWithError(w, http.StatusUnauthorized, securityErr.Error())
		return
	}

	token := r.Header.Get("ApplicationToken")
	nickname, dberr := QueryAssertToken(a.db, token)
	if dberr != nil {
		handleDBErrors(w, dberr)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"Nickname": nickname})
}
