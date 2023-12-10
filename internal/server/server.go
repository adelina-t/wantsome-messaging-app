package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	//"syscall"
)

// var shutdown os.Signal = syscall.SIGUSR1
var shutdown = os.Interrupt

func RunServer() {
	http.HandleFunc("/", home)
	http.HandleFunc("/ws", handleConnections)

	go handleMsg()

	server := &http.Server{Addr: ":8080"}

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
