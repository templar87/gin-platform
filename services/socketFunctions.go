package services

import (
	"net"
)

type SocketFunctions struct {
}

func (c SocketFunctions) Login(conn net.Conn, data interface{}) {
	source := data.(map[string]interface{})
	id := source["id"].(string)

	Socketconnections[id] = conn
}

func (c SocketFunctions) SendMessageToClient(data interface{}) {
	source := data.(map[string]interface{})

	to := source["to"].(string)
	toConn := Socketconnections[to]
	message := source["message"].(string)

	toConn.Write([]byte(message))
}
