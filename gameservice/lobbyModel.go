package main

import(
)

type Lobby struct{
	players [NUMPLAYERS]lobbyPlayer
	size int
}

type lobbyPlayer struct{
	exists bool
	connection *playerConnection
}

func (l *Lobby) addPlayer(connection *playerConnection){
	i := l.findFirstOpenLobbySpace()

	l.players[i].connection = connection
	l.players[i].exists = true
}

func (l* Lobby) numSpacesOpen() int{
	count := NUMPLAYERS

	for _, val := range l.players{
		if val.exists {
			count -= 1
		}
	}

	return count
}

//returns -1 if no space is open
func (l *Lobby) findFirstOpenLobbySpace() int{
	for i, val := range l.players{
		if val.exists {
			return i
		}
	}

	return -1
}
