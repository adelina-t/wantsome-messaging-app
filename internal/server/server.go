package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"wantsome.ro/messagingapp/pkg/utils"
)

var shutdown os.Signal = syscall.SIGKILL

var ServerLogger utils.Logger

func RunServer() {
	http.HandleFunc("/", home)
	http.HandleFunc("/rooms/", handleRoomConnection)
	http.HandleFunc("/chat/", handlePrivateChatConnection)

	go handlePrivateMessage()

	server := &http.Server{Addr: ":8080"}

	ServerLogger = utils.InitLogger()
	go ServerLogger.WriteToFileService()

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

	server.Shutdown(nil)
	if signal == shutdown {
		ServerLogger.DeInitLogger()
		os.Exit(1)
	}

}
