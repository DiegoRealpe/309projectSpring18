package main

import ( "fmt"
)

func ListenAndSend(g Game, connNumber int){
	for{
		readbyte, _ := g.reader[connNumber].Peek(1)
		id := int(readbyte[0])
		if id == 120 {
			bytemessage, err := g.reader[connNumber].Peek(17)
			if err != nil{
				fmt.Println("packet error: " + err.Error())
			}
			//parse the byte message
			//construct a message to broadcast to the clients
			//broadcast that message to the clients 
			
			g.reader[connNumber].Read(bytemessage)
		} else {
			g.reader[connNumber].ReadByte()
		}
	}
}
