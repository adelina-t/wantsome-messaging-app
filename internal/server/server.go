package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"wantsome.ro/messagingapp/internal/config"
)

var shutdown os.Signal = syscall.SIGTERM

func RunServer() {
	config := config.LoadConfig()

	serverAddress := fmt.Sprintf("%s:%s", config.URL, config.Port)

	http.HandleFunc("/", home)
	http.HandleFunc("/ws", handleConnections)

	go handleMsg()

	server := &http.Server{
		Addr: serverAddress,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Printf("Starting server on %s\n", serverAddress)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("error starting server: %s", err)
			stop <- shutdown
		}
	}()

	signal := <-stop
	log.Printf("Shutting down server ... ")

	m.Lock()
	for conn := range userConnections {
		conn.Close()
		delete(userConnections, conn)
	}
	m.Unlock()

	server.Shutdown(nil)
	if signal == shutdown {
		os.Exit(1)
	}

}
