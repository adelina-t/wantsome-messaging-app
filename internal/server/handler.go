package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"wantsome.ro/messagingapp/pkg/models"
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

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world from my server!")
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("got error upgrading connection %s\n", err)
		return
	}
	defer conn.Close()

	m.Lock()
	userConnections[conn] = ""
	m.Unlock()
	log.Printf("connected client!")

	for {
		var msg models.Message = models.Message{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("got error reading message %s\n", err)
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
		if msg.Message == "/list-users" {
			// Generate a list of user names
			userList := []string{}
			for _, userName := range userConnections {
				userList = append(userList, userName)
				log.Print(userList)
			}

			// Send the list back to the requesting user
			for conn, userName := range userConnections {
				if userName == msg.UserName {
					err := conn.WriteJSON(models.Message{Message: "Online users: " + strings.Join(userList, " ")})
					if err != nil {
						log.Printf("Error sending user list: %s", err)
					}
					break
				}
			}
		} else if msg.Recipient != "" {
			// send direct message
			for client, username := range userConnections {
				if username == msg.Recipient {
					err := client.WriteJSON(msg)
					if err != nil {
						log.Printf("Error sending direct message to %s: %s", msg.Recipient, err)
						client.Close()
						delete(userConnections, client)
					}
					break
				}
			}
		} else {
			// broadcast message
			for client, username := range userConnections {
				if username != msg.UserName {
					err := client.WriteJSON(msg)
					if err != nil {
						log.Printf("got error broadcating message to client %s", err)
						client.Close()
						delete(userConnections, client)
					}
				}
			}
		}
		m.Unlock()
	}
}
