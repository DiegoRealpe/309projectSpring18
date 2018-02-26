//testmain.go

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
CREATE TABLE IF NOT EXISTS users
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    age INT NOT NULL
)`

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

//TestEmptyTable test 1
func TestEmptyTable(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()
	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentUser(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()
	req, _ := http.NewRequest("GET", "/user/45", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "User not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'User not found'. Got '%s'", m["error"])
	}
}

func TestCreateUser(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()
	payload := []byte(`{"Nickname":"Knuckles"}`)
	req, _ := http.NewRequest("POST", "/client", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
	//var m map[string]interface{}
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
	req, _ := http.NewRequest("GET", "/client/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

}

func TestUpdateUser(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()
	addUsers(1)
	req, _ := http.NewRequest("GET", "/user/1", nil)
	response := executeRequest(req)
	var originalUser map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalUser)
	payload := []byte(`{"name":"test user - updated name","age":21}`)
	req, _ = http.NewRequest("PUT", "/user/1", bytes.NewBuffer(payload))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["id"] != originalUser["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalUser["id"], m["id"])
	}
	if m["name"] == originalUser["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalUser["name"], m["name"], m["name"])
	}
	if m["age"] == originalUser["age"] {
		t.Errorf("Expected the age to change from '%v' to '%v'. Got '%v'", originalUser["age"], m["age"], m["age"])
	}
}

func TestDeleteUser(t *testing.T) {
	testApp = App{}
	testApp.Initialize()

	clearTable()
	addUsers(1)
	req, _ := http.NewRequest("GET", "/user/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	req, _ = http.NewRequest("DELETE", "/user/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	req, _ = http.NewRequest("GET", "/user/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

/*********Helpers*********/

func addUsers(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		statement := fmt.Sprintf("INSERT INTO Clients(Nickname, GamesPlayed, GamesWon, GoalsScored, Active) VALUES('%s', %d, 0,0,0)", ("User " + strconv.Itoa(i+1)), ((i + 1) * 10))
		testApp.db.Exec(statement)
	}
}

func ensureTableExists() {
	if _, err := testApp.db.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	testApp.db.Exec("DELETE FROM Clients")
	testApp.db.Exec("ALTER TABLE Clients AUTO_INCREMENT = 1")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	fmt.Println("Executing")
	testApp.router.ServeHTTP(rr, req)

	return rr
}
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

/*db.Exec(`
CREATE TABLE Clients (
ID INT PRIMARY KEY,
Nickname VARCHAR(50) NOT NULL,
GamesPlayed INT NOT NULL,
GamesWon INT NOT NULL,
GoalsScored INT NOT NULL,
Online INT NOT NULL)`)*/
