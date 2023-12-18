package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"wantsome.ro/messagingapp/pkg/config"
)

var shutdown os.Signal = syscall.SIGUSR1

func RunServer() {
	// Load server configuration
	cfg := config.LoadConfig()

	http.HandleFunc("/", home)
	http.HandleFunc("/ws", handleConnections)

	go handleMsg()

	server := &http.Server{Addr: cfg.URL + ":" + cfg.Port}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Printf("Starting server on %s\n", server.Addr)
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
