package client

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"wantsome.ro/messagingapp/pkg/config"
	"wantsome.ro/messagingapp/pkg/models"
)

var (
	currentRoom string
	userName    string
)

func RunClient() {
	// Load server configuration
	cfg := config.LoadConfig()
	url := "ws://" + cfg.URL + ":" + cfg.Port + "/ws"

	if len(os.Args) < 2 {
		randId := rand.Intn(100)
		userName = fmt.Sprintf("Client%d", randId)
	} else {
		userName = os.Args[1]
	}

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("error dialing %s\n", err)
	}
	defer c.Close()

	log.Printf("Connected as user: %s", userName)

	done := make(chan bool)

	// Reading server messages
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Printf("error reading: %s\n", err)
				return
			}
			fmt.Printf("Got message: %s\n", message)
		}
	}()

	// Reading input from stdin
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input := scanner.Text()
			if strings.HasPrefix(input, "/") {
				handleCommands(c, userName, input)
			} else {
				message := models.Message{Message: input, UserName: userName, Room: currentRoom}
				sendMessage(c, message)
			}
		}
	}()
	<-done
	c.Close()
}

// Process messages from users
func handleCommands(c *websocket.Conn, user string, message string) {
	commands := strings.Split(message, " ")
	switch commands[0] {
	case "/quit":
		c.Close()
	case "/private":
		if len(commands) > 1 {
			rcp, txt := commands[1], strings.Join(commands[2:], " ")
			message := models.Message{Message: txt, UserName: user, Recipient: rcp}
			sendMessage(c, message)
		}
	case "/list-users":
		message := models.Message{Message: "/list-users", UserName: user}
		sendMessage(c, message)
	case "/list-rooms":
		message := models.Message{Message: "/list-rooms", UserName: user}
		sendMessage(c, message)
	case "/join":
		if len(commands) == 2 {
			room := commands[1]
			log.Printf("Joined room: %s\n", room)
			currentRoom = room
			message := models.Message{Message: "/join", UserName: user, Room: currentRoom}
			sendMessage(c, message)
		}
	case "/leave":
		if len(commands) == 2 {
			room := commands[1]
			log.Printf("Leaving room: %s\n", room)
			if currentRoom == room {
				currentRoom = ""
				log.Printf("Current room cleared")
			}
			message := models.Message{Message: "/leave", UserName: user, Room: room}
			sendMessage(c, message)
		} else {
			log.Printf("Leaving room: %s\n", currentRoom)
			message := models.Message{Message: "/leave", UserName: user, Room: currentRoom}
			sendMessage(c, message)
			currentRoom = ""
			log.Printf("Current room cleared")
		}
	case "/switch":
		if len(commands) == 2 {
			room := commands[1]
			log.Printf("Switched to room: %s\n", room)
			currentRoom = room
		}
	default:
		log.Printf("Sending universal command %s", message)
		message := models.Message{Message: message, UserName: user}
		sendMessage(c, message)
	}
}

func sendMessage(c *websocket.Conn, message models.Message) {
	err := c.WriteJSON(message)
	if err != nil {
		log.Printf("error writing %s\n", err)
		return
	}
}
