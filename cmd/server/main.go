package main

import (
	"log"

	"wantsome.ro/messagingapp/internal/server"
)

func main() {
	if err := server.RunServer(); err != nil {
		log.Fatalf("error running server: %s", err)
	}
}
