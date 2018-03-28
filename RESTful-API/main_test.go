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

var tableCreationQuery = `
CREATE TABLE Players (
ID INT PRIMARY KEY,
Nickname VARCHAR(50) NOT NULL,
GamesPlayed INT NOT NULL,
GamesWon INT NOT NULL,
GoalsScored INT NOT NULL,
Active INT NOT NULL)`

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
	req, _ := http.NewRequest("GET", "/player/45", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "" {
		t.Errorf("Expected the 'error' key of the response to be set to 'User not found'. Got '%s'", m["error"])
	}
}

func TestCreateUser(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()
	payload := []byte(`{"Nickname":"Knuckles"}`)
	req, _ := http.NewRequest("POST", "/player", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
	var jsonPlayer Player
	json.Unmarshal(response.Body.Bytes(), &jsonPlayer)

	if jsonPlayer.Nickname != "Knuckles" {
		t.Errorf("Expected user name to be 'Knuckles'. Got '%v', this user doesn't know de wey", jsonPlayer.Nickname)
	}
	if jsonPlayer.GamesPlayed != "0" {
		t.Errorf("Expected games won in new user to be to be '0'. Got '%v'", jsonPlayer.GamesPlayed)
	}
	if jsonPlayer.GoalsScored != "0" {
		t.Errorf("Expected goals scored to be '0'. Got '%v'", jsonPlayer.GoalsScored)
	}
}

func TestGetUser(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()
	addUsers(1)
	req, _ := http.NewRequest("GET", "/player/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

}

func TestUpdateUser(t *testing.T) {
	//Initialize routes and db
	testApp = App{}
	testApp.Initialize()
	clearTable()
	addUsers(1)

	//Get player that was just added
	req, _ := http.NewRequest("GET", "/player/1", nil)
	response := executeRequest(req)
	//Unmarshal the result
	var jsonPlayer Player
	json.Unmarshal(response.Body.Bytes(), &jsonPlayer)

	//Update Player
	payload := []byte(`{"Nickname":"newname","GamesPlayed":"21"}`)
	req, _ = http.NewRequest("PUT", "/player/1", bytes.NewBuffer(payload))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	//Modifying updated object
	jsonPlayer.Nickname = "newname"
	jsonPlayer.GamesPlayed = "21"

	//Get Resulting modified object
	var jsonPlayerR Player
	json.Unmarshal(response.Body.Bytes(), &jsonPlayerR)

	if jsonPlayer.ID != jsonPlayerR.ID {
		t.Errorf("Expected the id to remain the same (%v). Got %v", jsonPlayer.ID, jsonPlayerR.ID)
	}
	if jsonPlayer.GamesPlayed != jsonPlayerR.GamesPlayed {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", jsonPlayer.GamesPlayed, "21", jsonPlayerR.GamesPlayed)
	}
	if jsonPlayer.Nickname != jsonPlayerR.Nickname {
		t.Errorf("Expected the age to change from '%v' to '%v'. Got '%v'", jsonPlayer.Nickname, "newname", jsonPlayerR.Nickname)
	}
}

func TestDeleteUser(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()
	addUsers(1)
	req, _ := http.NewRequest("GET", "/player/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	req, _ = http.NewRequest("DELETE", "/player/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	req, _ = http.NewRequest("GET", "/player/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestRegisterUser(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()

	payload := []byte(`{"Nickname":"dumdum1"}`)
	req, _ := http.NewRequest("POST", "/player/register", bytes.NewBuffer(payload))
	req.Header.Set("FacebookToken", testUserToken)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	payload = []byte(`{"Nickname":"dumdum2"}`)
	req, _ = http.NewRequest("POST", "/player/register", bytes.NewBuffer(payload))
	req.Header.Set("FacebookToken", testUserToken2)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestLoginUser(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()

	payload := []byte(`{"Nickname":"User 1"}`)
	req, _ := http.NewRequest("POST", "/player/register", bytes.NewBuffer(payload))
	req.Header.Set("FacebookToken", testUserToken)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	req, _ = http.NewRequest("GET", "/player/1/login", nil)
	req.Header.Set("FacebookToken", testUserToken)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestCheckToken(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()

	payload := []byte(`{"Nickname":"User 1"}`)
	req, _ := http.NewRequest("POST", "/player/register", bytes.NewBuffer(payload))
	req.Header.Set("FacebookToken", testUserToken)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	testAppToken := m["ApplicationToken"]

	req, _ = http.NewRequest("GET", "/internal/checkApplicationToken", nil)
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
	req.Header.Set("FacebookToken", testUserToken)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	s, err := QueryAssertToken(testApp.db, "486074")
	if s != "dumdum1" {
		t.Errorf("returning nickname expected to be 'dumdum1'. Got '%s'", s)
	}
	if err != nil {
		t.Errorf(err.Error())
	}

}

//to get a new token login to facebook and get one from one of our test user
const testUserToken = "EAACqvTZC1964BAJe75cMt5tRyWcEoXo6fhFEjuaq8uXZAolqPSObU1AQ3LV59UiPMwW0ZC5dUgSCw9tin24GqMFZAMN67aST9NpZBtH5LDXvgWgUDEyhKMn7giZAG4H0UgmlRwSVhWRY8Qdagn87ZA8WdZAYvheA2vgXee4ZA7YCMY1unZBySRselJ7SZAwtCN0OqvTdkfzxZCweCNhrGnpLuGyrCFn683pbKkQbMkGAlmzixAZDZD"
const testUserToken2 = "EAACqvTZC1964BALDmZBZAONvSUYWBLq7Cvzxy3jNYmyX9GkvsZA7BJ8kxuh62Ekz4Fz0qJaZCAbeEWJcsCrR7Rtbev5zX8Ax5qRTaKd2ZAwYQHppZBz6ZCd8ZAR2LBJcZCE5HaGheeHGlO1eb80ipTp2OLBKxqphpk0GgKKqjO81yj8mAEvDBPkUCd4rMXnAKWbvA57YAS8pVKZAA3Hg5ZC5HZAeCjwMKOKuYQa0GHZCyZB0N9C5wZDZD"

func TestFBApiAccess(t *testing.T) {
	getFBUser(testUserToken)
}

/*********Test Helpers*********/

func addUsers(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		statement := fmt.Sprintf("INSERT INTO Players(Nickname, GamesPlayed, GamesWon, GoalsScored, Active) VALUES('%s', %d, 0,0,0)", ("User " + strconv.Itoa(i+1)), ((i + 1) * 10))
		testApp.db.Exec(statement)
	}
}

func ensureTableExists() {
	if _, err := testApp.db.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	testApp.db.Exec("DELETE FROM TokenTable")
	testApp.db.Exec("DELETE FROM FacebookData")
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
