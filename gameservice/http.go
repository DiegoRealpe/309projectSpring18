package main

import (
	"io"
	"net"
	"fmt"
	"net/http"
	"strconv"
	"encoding/json"
	"io/ioutil"
)

type clientConnection struct {
	connection net.Conn
	port int

	playerInfo connectionPlayerInfo
}

type connectionPlayerInfo struct {
	Id string `json:"id"`
	Username string `json:"nickname"`
	Gamesplayed string `json:"gamesplayed"`
	Gameswon string `json:"gameswon"`
	Goalsscored string `json:"goalsscored"`
	Rankwin string `json:"rankwin"`
	Rankscore string `json:"rankscore"`
	Lastavatar string `json:"lastavatar"`
}

type portHttpController struct {
	connPasser chan clientConnection
}


func makePortHttpController() portHttpController {
	connPasser := make(chan clientConnection)
	return portHttpController{connPasser: connPasser}
}


//when an http request is sent, send the requester a port and start listening on that port
func (portHttpController *portHttpController) handlePortRequested(w http.ResponseWriter, r *http.Request) {

	if numPortsAvailable() <= 0 {
		io.WriteString(w, "no ports avaliable, sorry fam")
		return
	}

	usedport := requestPort()

	stringport := strconv.Itoa(usedport)

	token, _ := strconv.Atoi(r.Header.Get("AppToken"))

	fmt.Println(token)

	apptoken := int64(token)

	io.WriteString(w, stringport)

	go func(apptoken int64) { //accept the first attempted connection on the port
		ln, _ := net.Listen("tcp", ":"+stringport)

		fmt.Println(apptoken)

		conn, _ := ln.Accept()

		ln.Close()// close connection so no new connections are accepted after player has quit

		connClient := clientConnection{
			connection: conn,
			port:				usedport,
		}

		connClient.playerInfo = checkTokenWithCrudService(apptoken)

		fmt.Println("nick:" + connClient.playerInfo.Username)

		portHttpController.connPasser <- connClient
	}(apptoken)

}


func checkTokenWithCrudService(internlToken int64) connectionPlayerInfo {
	Info := connectionPlayerInfo{}
	client := http.Client{}
	url := "http://proj-309-mg-6.cs.iastate.edu:8000/player/"
	strtoken:= strconv.Itoa(int(internlToken))
	url = url + strtoken
	fmt.Println(url)
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("ApplicationToken", strconv.Itoa(int(internlToken)))
	request.Header.Set("AppUser", "MG_6")
	request.Header.Set("AppSecret", "goingforthat#1bois")

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("crud request error")
		panic(err)
	}

	bodyByte, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.Body)
	err3 := json.Unmarshal(bodyByte, &Info)
	if err3 != nil {
		fmt.Println("unmarshal error")
		panic(err3)
	}

	return Info
}
