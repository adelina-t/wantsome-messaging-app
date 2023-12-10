package models

import "github.com/gorilla/websocket"

type Message struct {
	Content string
	From    string
	To      string // this set only for private messages
}

type Client struct {
	RoomID string
	Conn   *websocket.Conn
}

type ServerConfig struct {
	Address string
	Port    string
}

type ClientConfig struct {
	ClientID string
	RoomID   string
}
