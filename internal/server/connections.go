package server

import (
	"fmt"
	"log"
	"sync"

	"wantsome.ro/messagingapp/pkg/models"
)

var (
	m       sync.Mutex
	clients = map[string]*models.Client{}
)

func saveClient(id string, client *models.Client) {
	m.Lock()
	defer m.Unlock()

	clients[id] = client
	log.Printf("client(%s) connected to room(%s)", id, client.RoomID)
}

func deleteClient(id string) {
	client, exists := clients[id]
	if !exists {
		log.Printf("client(%s) is not connected", id)
		return
	}

	m.Lock()
	defer m.Unlock()

	if err := client.Conn.Close(); err != nil {
		log.Printf("error closing client(%s) ws connection: %s", id, err)
	}

	delete(clients, id)
	log.Printf("client(%s) disconnected", id)
}

func writeClientMessage(id string, message models.Message) error {
	client, ok := clients[id]
	if !ok {
		return fmt.Errorf("client(%s) is not connected", id)
	}

	m.Lock()
	defer m.Unlock()

	return client.Conn.WriteJSON(message)
}
