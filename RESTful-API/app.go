// app.go

package main

import (
	"database/sql"
	"encoding/json"
	"strconv"

	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//App struct where we have the db and router
type App struct {
	router *mux.Router
	db     *sql.DB
}

//Initialize the db and router
func (a *App) Initialize() {
	var err error
	a.db, err = sql.Open("mysql",
		"dbu309mg6:1XFA40wc@tcp(mysql.cs.iastate.edu:3306)/db309mg6")
	if err != nil {
		fmt.Println(err)
	}
	err = a.db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	a.router = mux.NewRouter()
	fmt.Println("Initialized")
}

//Run listen and serve
func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":8000", a.router))
}

func (a *App) getPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	p := Player{ID: string(id)}
	if err := p.getUser(a.db); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
