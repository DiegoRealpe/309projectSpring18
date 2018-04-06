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
		username : "∆∆∆™∆∆∆∆🐥🇺🇸語",
		connection: newPlayer.connection,
	}

	l.tellOtherPlayersYouJoined(newPlayer,out)
}

func (l *Lobby) sendExistingLobbyData(newPlayer *waitingPlayer) {
	l.sendAllExistingPlayers(newPlayer.connection.packetOut)
	l.sendAllChatMessagePackets(newPlayer.connection.packetOut)
}


func (l *Lobby) sendAllChatMessagePackets(to chan<- PacketOut){

	fmt.Println("messages are",l.messages)

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

	fmt.Println("repeating message",messageIn.message)


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

func (l *Lobby) playerNumberForConnectionID(id int) int{
	fmt.Println(l.players,l.size)
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

	fmt.Println("id, players",id,l.players)
	for i := 0; i < l.size ; i+= 1 {
		player := l.players[i]
		if player.connection.id != id{
			rtn[rtnIndex] = player.connection.id
			rtnIndex += 1
		}
	}

	fmt.Println("sending to",rtn)
	return rtn
}