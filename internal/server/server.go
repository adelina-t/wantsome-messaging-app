package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"wantsome.ro/messagingapp/internal/config"
)

var (
	shutdown os.Signal = syscall.SIGUSR1
)

func getServer() *http.Server {
	// get server config
	cfg := config.GetServerConfig()

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%s", cfg.Address, cfg.Port),
	}

	// enable verbose requests logging
	if ok, _ := strconv.ParseBool(os.Getenv("ENABLE_VERBOSE")); ok {
		server.Handler = logRequest(http.DefaultServeMux)
	}

	return server
}

func RunServer() error {
	// get server
	server := getServer()

	// register handlers
	http.HandleFunc("/", homeHttpHandler)
	http.HandleFunc("/ws", clientWsHandler)
	http.HandleFunc("/send-msg", sendMessageHttpHandler)
	http.HandleFunc("/list-users", listUsersHttpHandler)
	http.HandleFunc("/list-rooms", listRoomsHttpHandler)

	// add signal handler for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// start server into a goroutine
	go func() {
		log.Printf("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("error starting server: %s", err)
			stop <- shutdown
		}
	}()

	// graceful server shutdown
	signal := <-stop
	log.Printf("Shutting down server...")

	for clientID := range clients {
		deleteClient(clientID)
	}

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("error shutting down server: %s", err)
	}
	if signal == shutdown {
		return fmt.Errorf("server shutdown")
	}

	return nil
}
