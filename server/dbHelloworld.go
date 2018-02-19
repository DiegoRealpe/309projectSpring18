package main;

import(
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main(){
	fmt.Println("starting db test")

	db, err := sql.Open("mysql",
		"dbu309mg6:1XFA40wc@tcp(mysql.cs.iastate.edu:3306)/db309mg6")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()


	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	query1(db)

}

func query1(db *sql.DB){
	var rows , err = db.Query("select id, value from ryan_test")
	if err != nil {
		fmt.Println(err)
	}

	var id int
	var value string
	for rows.Next() {
		err := rows.Scan(&id, &value)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("name = ",id, "value = " , value)
}

