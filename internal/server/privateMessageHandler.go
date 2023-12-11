package server

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"wantsome.ro/messagingapp/pkg/models"
)

type SafePrivateChatList struct {
	privateChatList map[*websocket.Conn]string
	chatListMutex   sync.Mutex
	chatBroadcast   chan models.PrivateMessage
}

type UserList struct {
	UserList string `json:"userList"`
}

type RoomList struct {
	RoomList string `json:"roomList"`
}

func handlePrivateMessage() {
	loggerMethod := "handlePrivateMessage"
	loggerMessage := ""
	for {
		privMsg := <-privateChatList.chatBroadcast

		receivingUser := privMsg.ReceivingUser
		privateChatList.chatListMutex.Lock()
		receivingConnection := privateChatList.GetUserConnection(receivingUser)
		if receivingConnection != nil {
			err := receivingConnection.WriteJSON(privMsg)
			if err != nil {
				logMessage = fmt.Sprintf("got error broadcating message to client %s", err)
				ServerLogger.Log(logMessage, loggerMethod, models.LogLevel(0))
			}
		} else {
			loggerMessage = fmt.Sprintf("Could not find user %s", receivingUser)
			ServerLogger.Log(loggerMessage, loggerMethod, models.LogLevel(2))
		}
		privateChatList.chatListMutex.Unlock()
	}
}

func (safeChatList *SafePrivateChatList) GetUserConnection(desiredUser string) *websocket.Conn {
	for client, username := range safeChatList.privateChatList {
		if username == desiredUser {
			return client
		}
	}
	return nil
}

func listChatUsers(safePrivateChatList *SafePrivateChatList, conn *websocket.Conn) {
	logMethod := "listChatUsers"
	logMessage := ""
	userListString := ""
	safePrivateChatList.chatListMutex.Lock()
	for _, username := range safePrivateChatList.privateChatList {
		userListString += username + "\n"
	}
	userListStruct := UserList{
		UserList: userListString,
	}
	err := conn.WriteJSON(userListStruct)
	if err != nil {
		logMessage = fmt.Sprintf("Error sending the chat list %s", err)
		ServerLogger.Log(logMessage, logMethod, models.Error)
		conn.Close()
		delete(safePrivateChatList.privateChatList, conn)
	}
	safePrivateChatList.chatListMutex.Unlock()
}

func listRooms(conn *websocket.Conn, safeRooms *SafeTypeRoomList) {
	logMethod := "listChatUsers"
	logMessage := ""
	roomListString := ""

	safeRooms.m.Lock()
	for _, roomName := range safeRoomList.roomList {
		roomListString += roomName.Room.RoomName + "\n"
	}

	roomListStruct := RoomList{
		RoomList: roomListString,
	}

	err := conn.WriteJSON(roomListStruct)
	if err != nil {
		logMessage = fmt.Sprintf("Error sending the chat list %s", err)
		ServerLogger.Log(logMessage, logMethod, models.Error)
		conn.Close()
	}

	safeRooms.m.Unlock()

}
