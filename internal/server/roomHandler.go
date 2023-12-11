package server

import (
	"fmt"
	"slices"
	"sync"

	"github.com/gorilla/websocket"
	"wantsome.ro/messagingapp/pkg/models"
)

type SafeTypeRoomList struct {
	roomList []models.SafeRoom
	m        sync.Mutex
}

type RoomUserList struct {
	RoomUserList string `json:"roomUserList"`
}

func (safeRoomList *SafeTypeRoomList) AddNewRoom(newRoom models.SafeRoom) {
	safeRoomList.m.Lock()
	safeRoomList.roomList = append(safeRoomList.roomList, newRoom)
	safeRoomList.m.Unlock()
}

func (safeRoomList *SafeTypeRoomList) GetRoomIndex(roomName string) int {
	index := slices.IndexFunc(safeRoomList.roomList, func(localRoom models.SafeRoom) bool {
		return localRoom.Room.RoomName == roomName
	})
	return index
}

func handleRoomMessages(safeRoom *models.SafeRoom) {
	logMethod := "handleRoomMessages"
	logMessage := ""
	for {
		msg := <-safeRoom.Room.RoomBroadcast

		safeRoom.RoomMutex.Lock()
		ServerLogger.Log("WritingMessage", logMethod, models.LogLevel(0))
		for client, username := range safeRoom.Room.UserChatList {
			if username != msg.UserName {
				err := client.WriteJSON(msg)
				if err != nil {
					logMessage = fmt.Sprintf("got error broadcating message to client %s", err)
					ServerLogger.Log(logMessage, logMethod, models.LogLevel(0))
					client.Close()
					delete(safeRoom.Room.UserChatList, client)
				}
			}
		}
		safeRoom.RoomMutex.Unlock()
	}
}

func listRoomUsers(safeRoom models.SafeRoom, conn *websocket.Conn) {
	userList := ""
	safeRoom.RoomMutex.Lock()
	for _, username := range safeRoom.Room.UserChatList {
		userList += username + "\n"
	}

	roomUserListStruct := RoomUserList{
		RoomUserList: userList,
	}

	err := conn.WriteJSON(roomUserListStruct)
	if err != nil {
		fmt.Printf("got error listing users to chat %s", err)
		conn.Close()
		delete(safeRoom.Room.UserChatList, conn)
	}
	safeRoom.RoomMutex.Unlock()
}
