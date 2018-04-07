package main

import(
	"fmt"
)

type Lobby struct{
	players [NUMPLAYERS]lobbyPlayer
	size int

	messages []chatMessage
}

type lobbyPlayer struct{
	ready bool
	username string
	connection *playerConnection
}

type chatMessage struct{
	playerNumber int
	message string
}

func (l *Lobby) addPlayer(newPlayer *waitingPlayer, out chan<- PacketOut){
	i := l.size
	l.size += 1

	newPlayer.connection.packetOut <- PacketOut{ data:[]byte{222,byte(l.size-1)},size:2 }

	l.sendExistingLobbyData(newPlayer)

	l.players[i] = lobbyPlayer{
		ready: false,
		username : "∆∆∆™∆∆∆∆🐥🇺🇸語",
		connection: newPlayer.connection,
	}

	l.tellOtherPlayersYouJoined(newPlayer,out)
}

func (l *Lobby) sendExistingLobbyData(newPlayer *waitingPlayer) {
	l.sendAllExistingPlayers(newPlayer.connection.packetOut)
	l.sendAllChatMessagePackets(newPlayer.connection.packetOut)
	l.sendExistingPlayersReady(newPlayer.connection.packetOut)
}

func (l *Lobby) sendAllChatMessagePackets(to chan<- PacketOut){

	if debug {fmt.Println("messages are",l.messages)}

	for _ , m := range l.messages {
		message := packet203{
			playerNumber: byte(m.playerNumber),
			message: m.message,
		}

		packetOut := PacketOut{
			size: 402,
			data: message.toBytes(),
		}

		to <- packetOut
	}
}

func (l *Lobby) respondTo202(in *PacketIn, out chan<- PacketOut) {
	messageIn := ParseBytesTo202(in.data)
	playerNumber := l.playerNumberForConnectionID(in.connectionId)

	if debug {fmt.Println("repeating message",messageIn.message)}


	message := chatMessage{
		playerNumber: playerNumber,
		message: messageIn.message,
	}

	l.messages = append(l.messages, message)

	messageOut := packet203{
		playerNumber: byte(playerNumber),
		message: messageIn.message,

	}

	packetOut := PacketOut{
		size: 402,
		data: messageOut.toBytes(),
		targetIds: l.allConnectionIdsBut(in.connectionId),
	}

	out <- packetOut
}

func (l *Lobby) respondTo200(in *PacketIn, out chan<- PacketOut){
	playerNum := l.playerNumberForConnectionID(in.connectionId)
	if debug {fmt.Println("player",playerNum,"is ready")}
	l.players[playerNum].ready = true

	packet := packet204{numReady:1}
	packetOut := PacketOut{
		size: 2,
		data: packet.toBytes(),
		targetIds: l.allConnectionIdsBut(in.connectionId),
	}
	out <- packetOut

	if l.areAllPlayersReadyForTheGame() {
		fmt.Println("moving to game controller")
	}
}

func (l *Lobby) respondTo201(in *PacketIn, out chan<- PacketOut){
	playerNum := l.playerNumberForConnectionID(in.connectionId)
	if debug{fmt.Println("player",playerNum,"is unReady")}
	l.players[playerNum].ready = false

	packet := packet205{numUnready:1}
	packetOut := PacketOut{
		size: 2,
		data: packet.toBytes(),
		targetIds: l.allConnectionIdsBut(in.connectionId),
	}
	out <- packetOut

}

func (l *Lobby) playerNumberForConnectionID(id int) int{
	for i := 0; i < l.size; i += 1{
		if l.players[i].connection.id == id {
			return i
		}
	}
	return -1
}

func (l *Lobby) sendAllExistingPlayers (to chan<- PacketOut){
	for i := 0 ; i < l.size; i++{
		if i == l.size-1{
			continue //don't send players their own username
		}

		packet := packet206{
			playerNumber: i,
			username: l.players[i].username,
		}

		to <- PacketOut{
			size: 82,
			data: packet.toBytes(),
		}
	}
}

func (l *Lobby) tellOtherPlayersYouJoined(player *waitingPlayer, out chan<- PacketOut){
	packet := packet206{
		playerNumber: l.playerNumberForConnectionID(player.connection.id),
		username: "temp",
	}

	out <- PacketOut{
		size: 82,
		data: packet.toBytes(),
		targetIds: l.allConnectionIdsBut(player.connection.id),
	}
}

func (l *Lobby) allConnectionIdsBut(id int) []int{
	rtn := make([]int, l.size-1)
	rtnIndex := 0

	for i := 0; i < l.size ; i+= 1 {
		player := l.players[i]
		if player.connection.id != id{
			rtn[rtnIndex] = player.connection.id
			rtnIndex += 1
		}
	}

	return rtn
}


func (l *Lobby) areAllPlayersReadyForTheGame() bool{
	if l.size != NUMPLAYERS{
		return false
	}

	for _, val := range l.players{
		if !val.ready {
			return false
		}
	}

	return true
}

func (l *Lobby) sendExistingPlayersReady(out chan PacketOut) {
	numReady := 0

	//count players ready
	for i := 0; i < l.size ; i+= 1 {
		if l.players[i].ready {
			numReady += 1
		}
	}

	packet := packet204{numReady: byte(numReady)}
	packetOut := PacketOut{
		size: 2,
		data: packet.toBytes(),
	}

	out <- packetOut
}

func (lc *LobbyController) startGCFromLobby() {
	options := GameOptions{
		numPlayers: NUMPLAYERS,
		connectionIDToPlayerNumberMap: lc.getConnIDToPlayerNumberMap(),
	}
	for i, p := range lc.l.players {
		options.players[i] = p.connection
	}

	go runGameController(options, lc.packetIn, lc.packetOut)
}

func (lc *LobbyController) getConnIDToPlayerNumberMap() map[int]byte {
	idmap := make(map[int]byte)
	for i, p := range lc.l.players {
		idmap[p.connection.id] = byte(i)
	}
	return idmap
}

func (lc *LobbyController) startReadyTimer() bool {
	return true
}

func (l *Lobby) respondTo125(in *PacketIn, out chan<- PacketOut){
	fmt.Println("125, AHHHHHHHHHH")

	packetOut := PacketOut{
		size: 2,
		data: []byte{207,byte(in.connectionId)},
		targetIds: l.allConnectionIdsBut(in.connectionId),
	}
	out <-  packetOut

	for i := 0 ; i < l.size; i += 1 {
		l.players[i].connection.disconnect()
	}
}
func (lc *LobbyController) startGCFromLobby() {
	options := GameOptions{
		numPlayers: NUMPLAYERS,
		connectionIDToPlayerNumberMap: lc.getConnIDToPlayerNumberMap(),
	}
	for i, p := range lc.l.players {
		options.players[i] = p.connection
	}

	go runGameController(options, lc.packetIn, lc.packetOut)
}

func (lc *LobbyController) getConnIDToPlayerNumberMap() map[int]byte {
	idmap := make(map[int]byte)
	for i, p := range lc.l.players {
		idmap[p.connection.id] = byte(i)
	}
	return idmap
}

func (lc *LobbyController) startReadyTimer() bool {
	return true
}
