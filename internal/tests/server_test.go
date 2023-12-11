package tests

import (
	"log"
	"testing"

	"github.com/gorilla/websocket"
	"wantsome.ro/messagingapp/internal/server"
	"wantsome.ro/messagingapp/pkg/models"
)

func runTestServer() {
	server.RunServer("localhost", "3000")
}

func TestRoomChat_Connection_ListUser(t *testing.T) {
	go runTestServer()

	log.Printf("Starting the connection")
	testConn, _, err := websocket.DefaultDialer.Dial("ws://localhost:3000/rooms/room1", nil)
	log.Printf("Connection realised")
	if err != nil {
		log.Fatalf("error dialing %s\n", err)
	}
	defer testConn.Close()

	connectionMessage := models.GeneralMessage{
		Message:  "",
		UserName: "testDummy",
	}
	err = testConn.WriteJSON(connectionMessage)
	if err != nil {
		t.Fatalf("Failed to register user to room! %s", err)
	}

	connectionMessage = models.GeneralMessage{
		Message:  "!listUsers",
		UserName: "",
	}

	err = testConn.WriteJSON(connectionMessage)
	if err != nil {
		t.Fatalf("Failed to send list room users command! %s", err)
	}
	msg := server.RoomUserList{}
	err = testConn.ReadJSON(&msg)
	if err != nil {
		t.Fatalf("Failed to send list room users command! %s", err)
	}

	if msg.RoomUserList != "testDummy\n" {
		t.Fatalf("The test data is incorrect! %s", msg.RoomUserList)
	}
}
