package client

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"wantsome.ro/messagingapp/pkg/models"
)

func RunClient() {
	url := "ws://localhost:8080/ws"
	randId := rand.Intn(100)
	userName := fmt.Sprintf("Client%d", randId)

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
				message := models.Message{Message: input, UserName: userName}
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
	default:
		log.Printf("Unknown command")
	}
}

func sendMessage(c *websocket.Conn, message models.Message) {
	err := c.WriteJSON(message)
	if err != nil {
		log.Printf("error writing %s\n", err)
		return
	}
}
