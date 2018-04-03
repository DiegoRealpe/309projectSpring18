package main


type Lobby struct{
  members [NUMPLAYERS]*playerConnection
  isReady [NUMPLAYERS]int
  numMembers int
  connectionIDToPlayerNumberMap map[int]int
}
