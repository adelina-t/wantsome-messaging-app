package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"wantsome.ro/messagingapp/internal/config"
	"wantsome.ro/messagingapp/pkg/models"
)

func RunClientChat() error {
	// get client config
	clientCfg, err := config.GetClientConfig()
	if err != nil {
		return fmt.Errorf("error getting client config: %s", err)
	}

	// get server config
	serverCfg := config.GetServerConfig()
	serverUrl := fmt.Sprintf("ws://%s:%s/ws", serverCfg.Address, serverCfg.Port)

	// connect to server
	requestHeader := http.Header{
		"Client-ID": []string{clientCfg.ClientID},
		"Room-ID":   []string{clientCfg.RoomID},
	}
	conn, _, err := websocket.DefaultDialer.Dial(serverUrl, requestHeader)
	if err != nil {
		return fmt.Errorf("error dialing server: %s", err)
	}

	// add signal handler for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// print chat messages
	go func() {
		defer close(stop)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("error reading server meessage: %s", err)
				return
			}
			msg := models.Message{}
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Printf("error unmarshalling server message: %s", err)
				return
			}
			prefix := ""
			if msg.From != "" {
				// message from another client
				prefix = msg.From
				if msg.To != "" {
					// private message
					prefix += fmt.Sprintf(" -> %s (private)", msg.To)
				}
				prefix += ": "
			}
			fmt.Printf("%s%s\n", prefix, msg.Content)
		}
	}()

	// Disconnect client when the user presses Ctrl+C or connection is closed
	<-stop
	log.Printf("Exit client...")
	if err := conn.Close(); err != nil {
		log.Printf("error closing connection: %s", err)
	}

	return nil
}
