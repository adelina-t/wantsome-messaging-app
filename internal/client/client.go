package client

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"wantsome.ro/messagingapp/pkg/models"

	"bufio"
)

func RunClient() {
	url := "ws://localhost:8080/ws"
	//randId := rand.Intn(10)
	message := models.Message{
		Message:  "",
		UserName: "", //fmt.Sprintf("Client %d", randId),
	}

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("error dialing %s\n", err)
	}
	defer c.Close()

	//set username
	fmt.Print("Enter your username: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	message.UserName = scanner.Text()

	done := make(chan bool)

	// reading server messages
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Printf("error reading: %s\n", err)
				return
			}
			fmt.Printf("Got message from server: %s\n", message)
		}
	}()

	// writing messages to server
	go func() {
		scanner := bufio.NewScanner(os.Stdin)

		for {
			fmt.Printf("Enter a message: ")
			scanner.Scan()
			message.Message = scanner.Text()

			err := c.WriteJSON(message)
			if err != nil {
				log.Printf("error writing %s\n", err)
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()

	<-done
}
