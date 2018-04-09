package crudservice

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("Working")

	a := App{}
	a.Initialize()
	defer a.db.Close()

	a.Run()
	return
}
