package main

import (
	"fmt"
	"log"

	"github.com/gookit/config/v2"
	"wantsome.ro/messagingapp/internal/server"
)

func main() {

	err := config.LoadFiles("config.json")
	if err != nil {
		log.Printf("Could not load config %s", err)
		panic(err)
	}

	var serverHost string = fmt.Sprintf("%s", config.Data()["host"])
	var serverPort string = fmt.Sprintf("%s", config.Data()["port"])
	server.RunServer(serverHost, serverPort)
}
