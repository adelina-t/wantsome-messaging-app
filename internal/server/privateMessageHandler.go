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
