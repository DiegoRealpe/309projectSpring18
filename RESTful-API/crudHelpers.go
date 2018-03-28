package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

/*********Model Helpers*********/

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

//Helper functions that returns a future epoch time
//calculated from the time of the call plus a number of days given by parameter
func getExpiration(days int) int64 {
	epochDays := 86400 * days //epoch day lenght for every day in parameter
	return (time.Now().Unix() + int64(epochDays))
}

/*********Controller Helpers*********/

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
