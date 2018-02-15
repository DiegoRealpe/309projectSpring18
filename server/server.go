package main

import ( "fmt"
	"net"
	"net/http"
	"bufio"
	"strconv"
	"io"
)

var ports []int

func handler(w http.ResponseWriter, r *http.Request){

	if len(ports) <= 1 {
		io.WriteString(w, "no ports avaliable, sorry fam")
		return
	}

	usedport := ports[len(ports) - 1]

	ports = ports[:len(ports) - 1]

	stringport := strconv.Itoa(usedport)

	io.WriteString(w, stringport)


	go func() {
		ln, _ := net.Listen("tcp",":" + stringport)
		
		conn, _ := ln.Accept()
		
		writer := bufio.NewWriter(conn)
		writer.WriteString("welcome to port " + stringport + " :)\n")
		writer.Flush()
		conn.Close()

		ports = append(ports, usedport)
	}()

}



func main(){
	ports = []int{0, 3127, 8476, 1736, 5543, 9078}	

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":132", nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}
