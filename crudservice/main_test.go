package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var testApp App

//to get a new token login to facebook and get one from one of our test user

const tok = "EAACqvTZC1964BAL2gZCqzlzJNmvs0XDyFUnhMeGCqUB9afu8yFkmSgI5oLakRJWP8lZBq7z85pYz2SK4StQMZBFYk49cFmMglaBv9tuEQZBFTZCZBAN27RThzSezGOJYOieZCDn60fAKa01aWgdWd1xZCJtZBIlTciqdE8fdVomzTP9QZDZD"
const tok2 = "EAACqvTZC1964BAIrQcZAfscdEAebGdE5mWMheieBOObRMcZCAEjOXHCuCEqlttHpj8iZCPEcHqTQ7bfhJc8HWS9fmSrZCevqLwo2Un7wCoIeBx8ZCtpjKsRsC1vd9xoZA5cN3jVa9l1SKZButuOZBBh7M3j5FlfcCHFckEZC5nZBGOxg9xJGNFPCZA5A"
const tok3 = "EAACqvTZC1964BAGTnB9Iu2Rq5MXxiwMqhJh5V6FoTy0e6Gwh0ZCdjYlt29hhZAYm11mwXYBdGJDnFMDxFr1Sm6oWCQWlUWlyoA09bIj7rQuSYob3QdKqB4zy7lHrexBpecHFbYeLQAZBRX8EOWyugNAfMIVYZBbdQWRWWyrdTbC0vTZCXj8OaT"
const tok4 = "EAACqvTZC1964BAKjWuwNnziCBBDx5puhvU8yIlJyZBGMaJX6EzwzFb7lSZC4oo7w5O3c85QeFuZBYQn243fMgn9sJ2hNBlTH2ONGKBOXaZCXaESajcLc9U8RKOueNDRz18kCUThjpPzRBZCV1ZC1yISDTCg095lmK3HR7PGvzX6y57eZBuAlxYqX"
const tok5 = "EAACqvTZC1964BAHQpHTwBk0L9Xhq8tsliKUHVENyGO8tFTohDIVAtFBHGA7i5FVIboF6juuV8mog0JfBT9xZBm8l5zoifUFGC33OdFHKE99EEbQfayX40uUw0WS76ZC63uHQaTNlZAUUkJhZC8Y8yCewPIvMPuTpPPjGzF9cECZCh5kWlacgU1"
const tok6 = "EAACqvTZC1964BANFWVcg5mZAQRcZCJrfj1uLMz9uyxZA7BdEscEfyvDGU43pEYZCfrSgZAgthDJPofgEHlkUZAABtmDL6zj2zZAUC0s8MLJNzwrrZAwtCvsZCG3RBKDwCoOUuNvcb6kpGDgK6oBLjBlPptH1dEU8mJjc8gYOtRYgJVxu1NIbkMqYDJ"

var tokens = []string{tok, tok2, tok3, tok4, tok5, tok6}

var tableCreationQuery = `
CREATE TABLE Players1 (
ID INT PRIMARY KEY,
Nickname VARCHAR(50) NOT NULL,
GamesPlayed INT NOT NULL,
GamesWon INT NOT NULL,
GoalsScored INT NOT NULL)`

//Main Testing
func Testmain(t *testing.M) {
	testApp = App{}
	testApp.Initialize()
	testApp.Run()

	ensureTableExists()
	code := t.Run()
	clearTable()
	os.Exit(code)
}

func TestGetNonExistentUser(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()
	req, _ := http.NewRequest("GET", "/player/486074", nil)
	req.Header.Set("AppUser", "MG_6")
	req.Header.Set("AppSecret", "goingforthat#1bois")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Player Not Found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'User not found'. Got '%s'", m["error"])
	}
}

func TestFillUpTable(t *testing.T) {
	testApp = App{}
	testApp.Initialize()
	clearTable()

	for i, tokenAt := range tokens {
		payload := []byte(`{"Nickname":"Dummy#` + strconv.Itoa(i) + `"}`)
		req, _ := http.NewRequest("POST", "/player/register", bytes.NewBuffer(payload))
		req.Header.Set("FacebookToken", tokenAt)
		response := executeRequest(req)

		var m PlayerProfile
		json.Unmarshal(response.Body.Bytes(), &m)
		if m.Error != "" {
			t.Errorf(m.Error)
		}
		fmt.Println(m.Profile.Nickname)
		checkResponseCode(t, http.StatusCreated, response.Code)
	}
}

func TestGetUser(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()
	payload := []byte(`{"Nickname":"dumdum1"}`)
	req, _ := http.NewRequest("POST", "/player/register", bytes.NewBuffer(payload))
	req.Header.Set("FacebookToken", tok)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	req, _ = http.NewRequest("GET", "/player/486074", nil)
	req.Header.Set("AppUser", "MG_6")
	req.Header.Set("AppSecret", "goingforthat#1bois")
	response = executeRequest(req)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "" {
		t.Errorf(m["error"])
	}
	checkResponseCode(t, http.StatusOK, response.Code)

}

func TestUpdateUser(t *testing.T) {
	//Initialize routes and db
	testApp = App{}
	testApp.Initialize()
	clearTable()

	payload := []byte(`{"Nickname":"User 1"}`)
	req, _ := http.NewRequest("POST", "/player/register", bytes.NewBuffer(payload))
	req.Header.Set("FacebookToken", tok)
	response := executeRequest(req)

	var pro PlayerProfile
	json.Unmarshal(response.Body.Bytes(), &pro)
	if pro.Error != "" {
		t.Errorf(pro.Error)
	}
	checkResponseCode(t, http.StatusCreated, response.Code)

	//Get player that was just added
	req, _ = http.NewRequest("GET", "/player/486074", nil)
	req.Header.Set("AppUser", "MG_6")
	req.Header.Set("AppSecret", "goingforthat#1bois")
	response = executeRequest(req)
	//Unmarshal the result
	var jsonPlayer Player
	json.Unmarshal(response.Body.Bytes(), &jsonPlayer)

	//Update Player
	payload = []byte(`{"Nickname":"newname","GamesPlayed":"21"}`)
	req, _ = http.NewRequest("PUT", "/player/486074", bytes.NewBuffer(payload))
	req.Header.Set("AppUser", "MG_6")
	req.Header.Set("AppSecret", "goingforthat#1bois")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "" {
		t.Errorf(m["error"])
	}

	//Modifying updated object
	jsonPlayer.Nickname = "newname"
	jsonPlayer.GamesPlayed = "21"

	//Get Resulting modified object
	var jsonPlayerR Player
	json.Unmarshal(response.Body.Bytes(), &jsonPlayerR)

	if jsonPlayer.ID != jsonPlayerR.ID {
		t.Errorf("Expected the id to remain the same (%s). Got %s", jsonPlayer.ID, jsonPlayerR.ID)
	}
	if jsonPlayer.GamesPlayed != jsonPlayerR.GamesPlayed {
		t.Errorf("Expected the played games to change from '%v' to '%v'. Got '%v'", jsonPlayer.GamesPlayed, "21", jsonPlayerR.GamesPlayed)
	}
	if jsonPlayer.Nickname != jsonPlayerR.Nickname {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", jsonPlayer.Nickname, "newname", jsonPlayerR.Nickname)
	}
}

func TestDeleteUser(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()

	payload := []byte(`{"Nickname":"User 1"}`)
	req, _ := http.NewRequest("POST", "/player/register", bytes.NewBuffer(payload))
	req.Header.Set("FacebookToken", tok)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	req, _ = http.NewRequest("GET", "/player/486074", nil)
	req.Header.Set("AppUser", "MG_6")
	req.Header.Set("AppSecret", "goingforthat#1bois")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/player/486074", nil)
	req.Header.Set("AppUser", "MG_6")
	req.Header.Set("AppSecret", "goingforthat#1bois")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusAccepted, response.Code)

	req, _ = http.NewRequest("GET", "/player/486074", nil)
	req.Header.Set("AppUser", "MG_6")
	req.Header.Set("AppSecret", "goingforthat#1bois")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestRegisterUser(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()

	payload := []byte(`{"Nickname":"dumdum1"}`)
	req, _ := http.NewRequest("POST", "/player/register", bytes.NewBuffer(payload))
	req.Header.Set("FacebookToken", tok)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	payload = []byte(`{"Nickname":"dumdum2"}`)
	req, _ = http.NewRequest("POST", "/player/register", bytes.NewBuffer(payload))
	req.Header.Set("FacebookToken", tok2)
	response = executeRequest(req)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "" {
		t.Errorf(m["error"])
	}
	checkResponseCode(t, http.StatusCreated, response.Code)

}

func TestLoginUser(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()

	req0, reqerr := http.NewRequest("GET", "/player/login", nil)
	if reqerr != nil {
		t.Errorf(reqerr.Error())
	}
	req0.Header.Set("FacebookToken", tok)
	response0 := executeRequest(req0)
	var na map[string]string
	json.Unmarshal(response0.Body.Bytes(), &na)
	if na["error"] == "" {
		t.Errorf("Expected an error message and recieved none")
	}
	checkResponseCode(t, http.StatusNotFound, response0.Code)

	payload := []byte(`{"Nickname":"dumdum1"}`)
	req1, _ := http.NewRequest("POST", "/player/register", bytes.NewBuffer(payload))
	req1.Header.Set("FacebookToken", tok)
	response1 := executeRequest(req1)
	checkResponseCode(t, http.StatusCreated, response1.Code)

	req2, err := http.NewRequest("GET", "/player/login", nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	req2.Header.Set("FacebookToken", tok)
	response2 := executeRequest(req2)

	var m map[string]string
	json.Unmarshal(response2.Body.Bytes(), &m)
	if m["ApplicationToken"] == "" {
		t.Errorf("Expected an App Token and recieved none")
	}
	checkResponseCode(t, http.StatusAccepted, response2.Code)

}

func TestCheckToken(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()

	payload := []byte(`{"Nickname":"User 1"}`)
	req, _ := http.NewRequest("POST", "/player/register", bytes.NewBuffer(payload))
	req.Header.Set("FacebookToken", tok)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	testAppToken := m["ApplicationToken"]

	req, _ = http.NewRequest("GET", "/internal/checkApplicationToken", nil)
	req.Header.Set("AppUser", "MG_6")
	req.Header.Set("AppSecret", "goingforthat#1bois")
	req.Header.Set("ApplicationToken", testAppToken)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestTokenQuery(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()

	payload := []byte(`{"Nickname":"dumdum1"}`)
	req, _ := http.NewRequest("POST", "/player/register", bytes.NewBuffer(payload))
	req.Header.Set("FacebookToken", tok)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
	s, err := QueryAssertToken(testApp.db, "486074")
	if s != "1" {
		t.Errorf("returning ID expected to be '1'. Got '%s'", s)
	}
	if err != nil {
		t.Errorf(err.Error())
	}
}
func TestFBApiAccess(t *testing.T) {
	getFBUser(tok)
}

/*********Test Helpers*********/

func ensureTableExists() {
	if _, err := testApp.db.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	testApp.db.Exec("DELETE FROM Players")
	testApp.db.Exec("ALTER TABLE Players AUTO_INCREMENT = 1")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	testApp.router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
