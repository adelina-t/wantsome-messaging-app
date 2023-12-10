package server

import (
	"fmt"
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
	rooms           = make(map[string]*Room)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world from my server!")
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("got error upgrading connection %s\n", err)
		return
	}
	defer conn.Close()

	var currentRoom *Room

	m.Lock()
	userConnections[conn] = ""
	m.Unlock()
	fmt.Printf("connected client!")

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("got error reading message %s\n", err)
			m.Lock()
			delete(userConnections, conn)
			m.Unlock()
			return
		}

		m.Lock()
		if msg.Type == "login" {
			// Add user to the map of connections
			userConnections[conn] = msg.UserName
		}
		m.Unlock()

		if msg.Type == "join_room" {
			if currentRoom != nil {
				currentRoom.leaveRoom(conn)
			}
			currentRoom = getOrCreateRoom(msg.Room)
			currentRoom.joinRoom(conn)
		} else if currentRoom != nil {
			currentRoom.broadcastToRoom(msg)
		}
		// Handle request for user list
		if msg.Type == "list_users" {
			userList := make([]string, 0, len(userConnections))
			for _, username := range userConnections {
				if username != "" {
					userList = append(userList, username)
				}
			}

			userListMessage := models.Message{
				Type:    "user_list",
				Message: strings.Join(userList, "; "),
			}
			// Send back to the requester
			conn.WriteJSON(userListMessage)
		} else {
			m.Lock()
			userConnections[conn] = msg.UserName
			m.Unlock()
			broadcast <- msg
		}
	}
}

func handleMsg() {
	for {
		msg := <-broadcast

		m.Lock()
		// Iterate through all connected clients
		for client, username := range userConnections {
			// Check if the client is the intended recipient
			if username == msg.Recipient {
				err := client.WriteJSON(msg)
				if err != nil {
					fmt.Printf("got error sending private message to client %s\n", err)
					client.Close()
					delete(userConnections, client)
				}
				break
			}
		}
		m.Unlock()
	}
}

func getOrCreateRoom(roomName string) *Room {
	if room, ok := rooms[roomName]; ok {
		return room
	}
	newRoom := &Room{
		Name:    roomName,
		Members: make(map[*websocket.Conn]bool),
	}
	rooms[roomName] = newRoom
	return newRoom
}

func (r *Room) joinRoom(conn *websocket.Conn) {
	r.Members[conn] = true
}

func (r *Room) leaveRoom(conn *websocket.Conn) {
	delete(r.Members, conn)
}

func (r *Room) broadcastToRoom(msg models.Message) {
	for conn := range r.Members {
		conn.WriteJSON(msg)
	}
}
