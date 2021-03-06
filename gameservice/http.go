package main

import (
	"io"
	"net"
	"fmt"
	"net/http"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"bytes"
)

type clientConnection struct {
	connection net.Conn
	port int

	playerInfo connectionPlayerInfo
}

type connectionPlayerInfo struct {
	apptoken int
	Id string `json:"id"`
	Nickname string `json:"nickname"`
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

	token, _ := strconv.Atoi(r.Header.Get("AppToken"))

	apptoken := int64(token)

	connClient := clientConnection{}
	Error := 0
	connClient.playerInfo, Error = checkTokenWithCrudService(apptoken)

	if Error == 1 {
			fmt.Println("Invalid client")
			return;
	}

	usedport := requestPort()
	connClient.port = usedport

	stringport := strconv.Itoa(usedport)

	io.WriteString(w, stringport)

	go func(apptoken int64) { //accept the first attempted connection on the port
		ln, _ := net.Listen("tcp", ":"+stringport)

		conn, _ := ln.Accept()

		connClient.connection = conn

		ln.Close()// close connection so no new connections are accepted after player has quit

		fmt.Println("New connection from: " + connClient.playerInfo.Nickname)

		portHttpController.connPasser <- connClient
	}(apptoken)

}


func checkTokenWithCrudService(internlToken int64) (connectionPlayerInfo, int) {
	Info := connectionPlayerInfo{}
	Info.apptoken = int(internlToken)
	client := http.Client{}
	Error := 0
	url := "http://proj-309-mg-6.cs.iastate.edu:8000/player/"
	strtoken:= strconv.Itoa(int(internlToken))
	url = url + strtoken
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("ApplicationToken", strconv.Itoa(int(internlToken)))
	request.Header.Set("AppUser", "MG_6")
	request.Header.Set("AppSecret", "goingforthat#1bois")

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("crud request error")
		panic(err)
	}

	fmt.Println("Got response with status code", resp.StatusCode)

	if resp.StatusCode != 200{
		Error = 1
	}

	bodyByte, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.Body)
	err3 := json.Unmarshal(bodyByte, &Info)
	if err3 != nil {
		fmt.Println("unmarshal error")
		panic(err3)
	}

	return Info, Error
}

func (g *Game) sendEndgameDataToCrud(winningTeam int){
	url := "http://proj-309-mg-6.cs.iastate.edu:8000/player/"
	for n, p := range g.players {
		go func(p *playerConnection, n int) {
			playerurl := url + strconv.Itoa(p.playerInfo.apptoken)
			client := http.Client{}
			if n <= 1 && winningTeam == 0{
				newgameswon, _ := strconv.Atoi(p.playerInfo.Gameswon)
				p.playerInfo.Gameswon = strconv.Itoa(newgameswon + 1)
			}
			if n > 1 && winningTeam == 1{
				newgameswon, _ := strconv.Atoi(p.playerInfo.Gameswon)
				p.playerInfo.Gameswon = strconv.Itoa(newgameswon + 1)
			}
			bod, _ := json.Marshal(p.playerInfo)
			requestput, _ := http.NewRequest("PUT", playerurl, bytes.NewBuffer(bod))
			requestput.Header.Set("AppUser", "MG_6")
			requestput.Header.Set("AppSecret", "goingforthat#1bois")
			requestput.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(requestput)
			fmt.Println("Got response with code ", resp.StatusCode)
			if err != nil {
				panic(err)
			}
		}(p.connection, n)
	}
}
