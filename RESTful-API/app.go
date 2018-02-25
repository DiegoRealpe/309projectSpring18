// app.go

package main

import (
	"database/sql"

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
