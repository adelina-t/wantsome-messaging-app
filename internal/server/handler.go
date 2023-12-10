package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"wantsome.ro/messagingapp/pkg/models"

	"log"
	"time"
)

const (
	LogLevelInfo  = "INFO"
	LogLevelWarn  = "WARN"
	LogLevelError = "ERROR"
)

var (
	m               sync.Mutex
	userConnections = make(map[*websocket.Conn]string)
	broadcast       = make(chan models.Message)
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func logInfo(message string) {
	log.Printf("[%s] [%s] %s", time.Now().Format(time.RFC3339), LogLevelInfo, message)
}

func logWarn(message string) {
	log.Printf("[%s] [%s] %s", time.Now().Format(time.RFC3339), LogLevelWarn, message)
}

func logError(message string, err error) {
	log.Printf("[%s] [%s] %s: %v", time.Now().Format(time.RFC3339), LogLevelError, message, err)
}

func home(w http.ResponseWriter, r *http.Request) {
	logInfo("Handling home request")
	fmt.Fprintf(w, "Hello world from my server!")
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//fmt.Printf("got error upgrading connection %s\n", err)
		logError("Error upgrading connection", err)
		return
	}
	defer conn.Close()

	m.Lock()
	userConnections[conn] = ""
	m.Unlock()
	//fmt.Printf("connected client!")
	logInfo("Client connected!")

	for {
		var msg models.Message = models.Message{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			//fmt.Printf("got error reading message %s\n", err)
			logError("Error reading message", err)

			m.Lock()
			delete(userConnections, conn)
			m.Unlock()
			return
		}
		m.Lock()
		userConnections[conn] = msg.UserName
		m.Unlock()
		broadcast <- msg
	}
}

func handleMsg() {
	for {
		msg := <-broadcast

		m.Lock()
		for client, username := range userConnections {
			if username != msg.UserName {
				err := client.WriteJSON(msg)
				if err != nil {
					//fmt.Printf("got error broadcating message to client %s", err)
					logError("Error broadcasting message to client", err)

					client.Close()
					delete(userConnections, client)
				}
			}
		}
		m.Unlock()
	}
}
