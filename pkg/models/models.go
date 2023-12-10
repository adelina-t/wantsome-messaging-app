package models

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ChatUser struct {
	Username   string
	Connection *websocket.Conn
}
type GeneralMessage struct {
	Message   string
	UserName  string
	TimeStamp time.Time
}

type PrivateMessage struct {
	Message       string
	SendingUser   string
	ReceivingUser string
	TimeStamp     time.Time
}

type PrivateChat struct {
	ChatBroadcast chan string
	ReceivingUser ChatUser
	TimeStamp     time.Time
}

type SafePrivateChat struct {
	PrivateChat      PrivateChat
	PrivateChatMutex *sync.Mutex
}

type Room struct {
	RoomName      string
	UserChatList  map[*websocket.Conn]string
	Message       []GeneralMessage
	RoomBroadcast chan GeneralMessage
}

type SafeRoom struct {
	Room      Room
	RoomMutex *sync.Mutex
}

type LogLevel int

const (
	Error LogLevel = iota
	Warning
	Verbose
)

type LoggerMessage struct {
	Message string
	Method  string
	Level   LogLevel
}
