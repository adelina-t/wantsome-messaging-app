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
	rooms           = make(map[string]map[*websocket.Conn]bool)
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
		if strings.HasPrefix(msg.Message, "/") {
			processCommands(msg)
		} else if msg.Recipient != "" {
			sendDirectMessage(msg)
		} else if msg.Room != "" {
			sendMessageToRoom(msg)
		} else {
			broadcastMessage(msg)
		}
		m.Unlock()
	}
}

func processCommands(msg models.Message) {
	switch msg.Message {
	case "/list-users":
		listUsers(msg)
	case "/list-rooms":
		listRooms(msg)
	case "/join":
		addUserToRoom(msg)
	case "/leave":
		removeUserFromRoom(msg)
	case "/broadcast":
		broadcastMessage(msg)
	default:
		log.Printf("Unknown command")
	}
}

func listUsers(msg models.Message) {
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
}

func listRooms(msg models.Message) {
	var roomNames []string
	for conn, userName := range userConnections {
		if userName == msg.UserName {
			for room, clients := range rooms {
				for client, _ := range clients {
					if client == conn {
						roomNames = append(roomNames, room)
					}
				}
			}
			err := conn.WriteJSON(models.Message{Message: "User " + msg.UserName + " is part of the following rooms: " + strings.Join(roomNames, " ")})
			if err != nil {
				log.Printf("Error sending user list: %s", err)
			}
			break
		}
	}

}

func addUserToRoom(msg models.Message) {
	roomName := msg.Room
	for conn, userName := range userConnections {
		if userName == msg.UserName {
			if _, ok := rooms[roomName]; !ok {
				rooms[roomName] = make(map[*websocket.Conn]bool)
			}
			rooms[roomName][conn] = true
			break
		}
	}
}

func removeUserFromRoom(msg models.Message) {
	for conn, userName := range userConnections {
		if userName == msg.UserName {
			delete(rooms[msg.Room], conn)
			break
		}
	}
}

func sendMessageToRoom(msg models.Message) {
	for room, clients := range rooms {
		if room == msg.Room {
			for client, _ := range clients {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("got error broadcating message to client %s", err)
					client.Close()
					delete(userConnections, client)
				}
			}
		}
	}
}

func broadcastMessage(msg models.Message) {
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

func sendDirectMessage(msg models.Message) {
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
}
