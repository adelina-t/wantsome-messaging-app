package config

import (
	"fmt"
	"os"

	"wantsome.ro/messagingapp/pkg/models"
)

const (
	defaultServerAddr = "localhost"
	defaultServerPort = "8080"

	defaultClientRoomID = "default"
)

func GetServerConfig() *models.ServerConfig {
	cfg := models.ServerConfig{
		Address: os.Getenv("SERVER_ADDR"),
		Port:    os.Getenv("SERVER_PORT"),
	}

	if cfg.Address == "" {
		cfg.Address = defaultServerAddr
	}
	if cfg.Port == "" {
		cfg.Port = defaultServerPort
	}

	return &cfg
}

func GetClientConfig() (*models.ClientConfig, error) {
	cfg := models.ClientConfig{
		ClientID: os.Getenv("CLIENT_ID"),
		RoomID:   os.Getenv("CLIENT_ROOM_ID"),
	}

	if cfg.ClientID == "" {
		return nil, fmt.Errorf("CLIENT_ID env variable is empty")
	}
	if cfg.RoomID == "" {
		cfg.RoomID = defaultClientRoomID
	}

	return &cfg, nil
}
