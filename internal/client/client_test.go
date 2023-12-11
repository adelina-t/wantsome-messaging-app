package client

import (
	"testing"

	"wantsome.ro/messagingapp/pkg/models"
)

func TestCreateLoginMessage(t *testing.T) {
	username := "testuser"
	expected := models.Message{
		Type:     "login",
		UserName: username,
	}
	result := createLoginMessage(username)
	if result != expected {
		t.Errorf("Expected %#v, got %#v", expected, result)
	}
}

func TestCreateJoinRoomMessage(t *testing.T) {
	username := "testuser"
	room := "testroom"
	expected := models.Message{
		Type:     "join_room",
		UserName: username,
		Room:     room,
	}
	result := createJoinRoomMessage(username, room)
	if result != expected {
		t.Errorf("Expected %#v, got %#v", expected, result)
	}
}

func TestCreateRoomMessage(t *testing.T) {
	username := "testuser"
	message := "testmessage"
	room := "testroom"
	recipient := "testrecipient"
	expected := models.Message{
		UserName:  username,
		Message:   message,
		Room:      room,
		Recipient: recipient,
	}
	result := createRoomMessage(username, message, room, recipient)
	if result != expected {
		t.Errorf("Expected %#v, got %#v", expected, result)
	}
}
