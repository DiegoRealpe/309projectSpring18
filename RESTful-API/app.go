// app.go

package main

import (
	"database/sql"
	"encoding/json"

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
	a.initializeRoutes()
	fmt.Println("Initialized")
}

//Run listen and serve
func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":8000", a.router))
}

func (a *App) initializeRoutes() {
	a.router.HandleFunc("/player/{ID}", a.getPlayer).Methods("GET")
	a.router.HandleFunc("/player", a.createPlayer).Methods("POST") //No mux params, credentials in request body
	a.router.HandleFunc("/player/{ID}", a.deletePlayer).Methods("DELETE")
	/*a.router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET")
	a.router.HandleFunc("/user/{id:[0-9]+}", a.updateUser).Methods("PUT")*/
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
