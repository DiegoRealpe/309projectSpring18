package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("Working")

	a := App{}
	a.Initialize()
	defer a.db.Close()

	return
	/*fmt.Println("Server")

	a.router.HandleFunc("/player", GetAllPlayers).Methods("GET")
	a.router.HandleFunc("/player/{id}", GetPlayer).Methods("GET")
	a.router.HandleFunc("/player/{id}", CreatePlayer).Methods("POST")
	a.router.HandleFunc("/player/{id}", DeletePlayer).Methods("DELETE")

	a.Run()*/
}
