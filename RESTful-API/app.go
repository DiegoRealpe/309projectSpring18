// app.go

package main

import (
	"database/sql"
	"encoding/json"
	"errors"
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
	a.initializeRoutes()
	fmt.Println("Initialized")
}

//Run listen and serve
func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":8000", a.router))
}

func (a *App) initializeRoutes() {
	a.router.HandleFunc("/client/{ID}", a.getPlayer).Methods("GET")
	a.router.HandleFunc("/user", a.createPlayer).Methods("POST") //No mux params, credentials in request body
	/*a.router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET")
	a.router.HandleFunc("/user/{id:[0-9]+}", a.updateUser).Methods("PUT")
	a.router.HandleFunc("/user/{id:[0-9]+}", a.deleteUser).Methods("DELETE")*/
}

/*********Routes*********/

func (a *App) getPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	p := Player{ID: strconv.Itoa(id)}
	if err := p.QuerySearchPlayer(a.db); err != nil {
		switch err {
		case errors.New("Empty"):
			respondWithError(w, http.StatusBadRequest, "Empty")
		case errors.New("Query Error"):
			respondWithError(w, http.StatusBadRequest, "Bad Query")
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var p Player
	decoder := json.NewDecoder(r.Body) //Passing credentials through http request body
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	err := p.QueryCreatePlayer(a.db)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, p)
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
