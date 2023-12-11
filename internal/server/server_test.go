package server

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"wantsome.ro/messagingapp/pkg/models"
)

func TestRunServer(t *testing.T) {
	// Start the server
	go RunServer()

	// Wait for the server to start
	time.Sleep(2 * time.Second)

	// Make a GET request to the server's home endpoint
	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Fatalf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	log.Printf("Server is running on http://localhost:8080")

	// Make a WebSocket connection to the server
	u := "ws://localhost:8080/ws"
	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("failed to connect to WebSocket: %v", err)
	}
	defer c.Close()

	log.Printf("Client connected!")

	// Test sending a login message
	loginMsg := models.Message{
		Type:     "login",
		UserName: "testUser",
	}
	err = c.WriteJSON(loginMsg)
	if err != nil {
		t.Errorf("error sending login message: %v", err)
	}

	// Test Joining a Room
	joinRoomMsg := models.Message{
		Type:     "join_room",
		UserName: "testUser",
		Room:     "testRoom",
	}
	if err := c.WriteJSON(joinRoomMsg); err != nil {
		t.Errorf("error sending join room message: %v", err)
	}

	// Test Sending a Message in a Room
	roomMessage := models.Message{
		Type:     "message",
		UserName: "testUser",
		Room:     "testRoom",
		Message:  "Hello, Room!",
	}
	if err := c.WriteJSON(roomMessage); err != nil {
		t.Errorf("error sending message in room: %v", err)
	}

	// Test Sending a Message to show users
	userListMessage := models.Message{
		Type:     "user_list",
		UserName: "testUser",
		Room:     "",
		Message:  "",
	}
	if err := c.WriteJSON(userListMessage); err != nil {
		t.Fatalf("failed to send message: %v", err)
	}
}
