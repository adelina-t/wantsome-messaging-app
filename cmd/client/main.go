package main

import (
	"log"

	"wantsome.ro/messagingapp/internal/client"
)

func main() {
	if err := client.RunClientChat(); err != nil {
		log.Fatalf("error running server: %s", err)
	}
}
