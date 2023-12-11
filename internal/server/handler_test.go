package server

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestGetOrCreateRoom(t *testing.T) {
	// Test case 1: Room already exists
	roomName := "testRoom"
	rooms[roomName] = &Room{
		Name:    roomName,
		Members: make(map[*websocket.Conn]bool),
	}

	room := getOrCreateRoom(roomName)
	if room == nil {
		t.Errorf("getOrCreateRoom(%s) returned nil, expected room", roomName)
	}

	// Test case 2: Room doesn't exist
	roomName = "newRoom"
	room = getOrCreateRoom(roomName)
	if room == nil {
		t.Errorf("getOrCreateRoom(%s) returned nil, expected room", roomName)
	}

	// Verify that the room was created
	if _, ok := rooms[roomName]; !ok {
		t.Errorf("getOrCreateRoom(%s) did not create the room", roomName)
	}
}
