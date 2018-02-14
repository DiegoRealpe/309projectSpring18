package main

import (
	"fmt"
	"net/http"
)

func main() {

}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type Noob struct {

}

func (n Noob ) ServeHTTP(w http.ResponseWriter, r *http.Request){
	
}