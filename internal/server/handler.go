package server

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"wantsome.ro/messagingapp/pkg/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var safeRoomList SafeTypeRoomList = SafeTypeRoomList{
	roomList: make([]models.SafeRoom, 0),
	m:        sync.Mutex{},
}

var privateChatList SafePrivateChatList = SafePrivateChatList{
	privateChatList: make(map[*websocket.Conn]string),
	chatListMutex:   sync.Mutex{},
	chatBroadcast:   make(chan models.PrivateMessage),
}

var logMessage string

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world from my server!")
}

func handleRoomConnection(w http.ResponseWriter, r *http.Request) {
	loggerMethod := "handleRoomConnection"
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("got error upgrading connection %s\n", err)
		return
	}
	defer conn.Close()

	roomName := strings.TrimPrefix(r.URL.Path, "/rooms/")
	if roomName == "" {
		roomName = "General"
	}

	var connRoom models.SafeRoom
	roomIndex := safeRoomList.GetRoomIndex(roomName)
	if roomIndex == -1 {
		var newRoom models.Room = models.Room{
			RoomName:      roomName,
			UserChatList:  make(map[*websocket.Conn]string),
			Message:       make([]models.GeneralMessage, 0),
			RoomBroadcast: make(chan models.GeneralMessage),
		}
		var newSafeRoom models.SafeRoom = models.SafeRoom{
			Room:      newRoom,
			RoomMutex: &sync.Mutex{},
		}
		safeRoomList.AddNewRoom(newSafeRoom)
		connRoom = newSafeRoom
		logMessage = fmt.Sprintf("Starting broadcasting for room %s", connRoom.Room.RoomName)
		ServerLogger.Log(logMessage, loggerMethod, models.LogLevel(2))

		go handleRoomMessages(&connRoom)
	} else {
		connRoom = safeRoomList.roomList[roomIndex]
	}

	if connRoom.Room.UserChatList[conn] != "" {
		logMessage = fmt.Sprintf("User %s is already connected!", connRoom.Room.UserChatList[conn])
		ServerLogger.Log(logMessage, loggerMethod, models.LogLevel(2))
		return
	}

	for {
		var msg models.GeneralMessage = models.GeneralMessage{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			logMessage = fmt.Sprintf("got error reading message %s", err)
			ServerLogger.Log(logMessage, loggerMethod, models.LogLevel(0))
			delete(connRoom.Room.UserChatList, conn)
			return
		}
		if msg.UserName != "" {
			if isUserAlreadyConnected(connRoom.Room.UserChatList, conn) {
				logMessage = fmt.Sprintf("User %s already in chat", connRoom.Room.UserChatList[conn])
				ServerLogger.Log(logMessage, loggerMethod, models.LogLevel(2))
			} else {
				logMessage = fmt.Sprintf("Setting the user %s", msg.UserName)
				ServerLogger.Log(logMessage, loggerMethod, models.LogLevel(2))
				connRoom.Room.UserChatList[conn] = msg.UserName
			}
		} else {
			msg.UserName = connRoom.Room.UserChatList[conn]

			if len(msg.Message) == 0 {
				logMessage = fmt.Sprintf("Empty message! User: %s", msg.UserName)
				ServerLogger.Log(logMessage, loggerMethod, models.LogLevel(2))
			}
			if !isCommand(msg.Message) {
				connRoom.Room.RoomBroadcast <- msg
			} else {
				listRoomUsers(connRoom, conn)
			}
		}

	}
}

func handlePrivateChatConnection(w http.ResponseWriter, r *http.Request) {
	loggerMethod := "handlePrivateChatConnection"
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("got error upgrading connection %s\n", err)
		return
	}
	defer conn.Close()

	if privateChatList.privateChatList[conn] != "" {
		logMessage = fmt.Sprintf("User %s is already connected!", privateChatList.privateChatList[conn])
		ServerLogger.Log(logMessage, loggerMethod, models.LogLevel(2))
		return
	}

	for {
		var msg models.PrivateMessage = models.PrivateMessage{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("got error reading message %s\n", err)
			//TODO: Make a function which removes the connection
			delete(privateChatList.privateChatList, conn)
			return
		}

		if msg.SendingUser != "" {
			if isUserAlreadyConnected(privateChatList.privateChatList, conn) {
				logMessage = fmt.Sprintf("User %s already in chat", privateChatList.privateChatList[conn])
				ServerLogger.Log(logMessage, loggerMethod, models.LogLevel(2))

			} else {
				logMessage = fmt.Sprintf("Setting the user %s", msg.SendingUser)
				ServerLogger.Log(logMessage, loggerMethod, models.LogLevel(2))
				privateChatList.privateChatList[conn] = msg.SendingUser
			}
		}
		if msg.ReceivingUser != "" && msg.Message != "" {
			privateChatList.chatBroadcast <- msg
		}
	}

}

func isUserAlreadyConnected(userConnections map[*websocket.Conn]string, conn *websocket.Conn) bool {
	if userConnections[conn] != "" {
		return true
	} else {
		return false
	}
}

func isCommand(message string) bool {
	return message[0] == '!'
}
