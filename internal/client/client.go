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
	randId := rand.Intn(10)
	userName := fmt.Sprintf("Client%d", randId)

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("error dialing %s\n", err)
	}
	defer c.Close()

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
			commands := strings.Split(input, " ")
			switch commands[0] {
			case "/quit":
				close(done)
				return
			case "/private":
				if len(commands) > 1 {
					message := models.Message{
						Message:   strings.Join(commands[1:], " "),
						UserName:  userName,
						Recipient: commands[1],
					}
					err := c.WriteJSON(message)
					if err != nil {
						log.Printf("error writing %s\n", err)
						return
					}
				}
			default:
				message := models.Message{Message: input, UserName: userName}
				err := c.WriteJSON(message)
				if err != nil {
					log.Printf("error writing %s\n", err)
					return
				}
			}
		}
	}()

	<-done
	c.Close()
}
