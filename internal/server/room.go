package server

import "github.com/gorilla/websocket"

type Room struct {
	Name    string
	Members map[*websocket.Conn]bool
}
