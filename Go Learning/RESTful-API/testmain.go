//testmain.go

package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testApp App

func testmain(t *testing.M) {
	testApp.Initialize()
	testApp.Run()

	code := t.Run()

	os.Exit(code)
}

//TestEmptyTable test 1
func TestEmptyTable(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func clearTable() {
	testApp.db.Exec("DELETE FROM users")
	testApp.db.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
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

/*db.Exec(`
CREATE TABLE Clients (
ID INT PRIMARY KEY,
Nickname VARCHAR(50) NOT NULL,
GamesPlayed INT NOT NULL,
GamesWon INT NOT NULL,
GoalsScored INT NOT NULL,
Online INT NOT NULL)`)*/
