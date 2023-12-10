package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"wantsome.ro/messagingapp/pkg/models"
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("RemoteAddr - %s", r.RemoteAddr)
		log.Printf("Proto - %s", r.Proto)
		log.Printf("Host - %s", r.Host)
		log.Printf("Method - %s", r.Method)
		log.Printf("RequestURI - %s", r.RequestURI)
		log.Printf("URL - %s", r.URL)
		log.Printf("Body - %s", r.Body)
		log.Print("Headers:")
		for k, v := range r.Header {
			log.Printf("    %s: %s", k, v)
		}
		handler.ServeHTTP(w, r)
	})
}

func homeHttpHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprintf(w, "Welcome to wantsome.ro messaging app!\n"); err != nil {
		log.Printf("error writing home http response: %s", err)
	}
}

func getClientIDHeader(w http.ResponseWriter, r *http.Request) (string, error) {
	clientID := r.Header.Get("Client-ID")
	if clientID == "" {
		return "", fmt.Errorf("Client-ID header is missing")
	}
	return clientID, nil
}

func getRoomIDHeader(w http.ResponseWriter, r *http.Request) (string, error) {
	roomID := r.Header.Get("Room-ID")
	if roomID == "" {
		return "", fmt.Errorf("Room-ID header is missing")
	}
	return roomID, nil
}

func clientWsHandler(w http.ResponseWriter, r *http.Request) {
	// upgrade the connection to websocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading connection to WS: %s", err)
		return
	}

	// validate "Client-ID" header
	clientID, err := getClientIDHeader(w, r)
	if err != nil {
		if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error())); err != nil {
			log.Printf("error writing client(%s) close message: %s", clientID, err)
		}
		if err := conn.Close(); err != nil {
			log.Printf("error closing client(%s) connection: %s", clientID, err)
		}
		return
	}
	if _, ok := clients[clientID]; ok {
		closeMessage := fmt.Sprintf("client(%s) already connected", clientID)
		if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, closeMessage)); err != nil {
			log.Printf("error writing client(%s) close message: %s", clientID, err)
		}
		if err := conn.Close(); err != nil {
			log.Printf("error closing client(%s) connection: %s", clientID, err)
		}
		return
	}

	// validate "Room-ID" header
	roomID, err := getRoomIDHeader(w, r)
	if err != nil {
		if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error())); err != nil {
			log.Printf("error writing client(%s) close message: %s", clientID, err)
		}
		if err := conn.Close(); err != nil {
			log.Printf("error closing client(%s) connection: %s", clientID, err)
		}
		return
	}

	// save the client connection and delete it when the function returns
	client := models.Client{
		RoomID: roomID,
		Conn:   conn,
	}
	saveClient(clientID, &client)
	defer deleteClient(clientID)

	// send welcome message to the client
	msg := models.Message{
		Content: fmt.Sprintf("Hello %s! Welcome to chat room %s!", clientID, roomID),
	}
	if err := writeClientMessage(clientID, msg); err != nil {
		log.Printf("error writing client(%s) message: %s", clientID, err)
		return
	}

	// loop until the client closes connection
	for {
		_, msgBytes, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("error reading client(%s) message: %s", clientID, err)
			return
		}

		msg := models.Message{}
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("error unmarshalling client(%s) message: %s", clientID, err)
			return
		}
	}
}

func sendMessageHttpHandler(w http.ResponseWriter, r *http.Request) {
	// validate "Client-ID" header
	clientID, err := getClientIDHeader(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := fmt.Fprintf(w, "error: %s\n", err); err != nil {
			log.Printf("error writing client http response: %s", err)
		}
		return
	}

	// validate that the reqClient is connected
	reqClient, connected := clients[clientID]
	if !connected {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := fmt.Fprintf(w, "error: client(%s) not connected\n", clientID); err != nil {
			log.Printf("error writing client http response: %s", err)
		}
		return
	}

	// validate proper HTTP method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := fmt.Fprintf(w, "error: invalid HTTP method %s, only POST is allowed\n", r.Method); err != nil {
			log.Printf("error writing client http response: %s", err)
		}
		return
	}

	// read HTTP body from request
	msgBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := fmt.Fprintf(w, "error reading body content %s\n", err); err != nil {
			log.Printf("error writing client http response: %s", err)
		}
		return
	}
	msg := models.Message{
		From:    clientID,
		Content: string(msgBytes),
	}

	// send client message
	var logMessage string

	toClientID := r.Header.Get("To-Client-ID")
	if toClientID != "" {
		// send private message
		msg.To = toClientID
		logMessage = fmt.Sprintf("sending private message from client(%s) to client(%s)", clientID, toClientID)
		if _, ok := clients[toClientID]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := fmt.Fprintf(w, "error: private message destination client(%s) not connected\n", toClientID); err != nil {
				log.Printf("error writing client http response: %s", err)
			}
			return
		}
		if err := writeClientMessage(clientID, msg); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := fmt.Fprintf(w, "error: %s\n", err); err != nil {
				log.Printf("error writing client http response: %s", err)
			}
			return
		}
		if err := writeClientMessage(toClientID, msg); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := fmt.Fprintf(w, "error: %s\n", err); err != nil {
				log.Printf("error writing client http response: %s", err)
			}
			return
		}
	} else {
		// broadcast message to all clients in the room
		logMessage = fmt.Sprintf("broadcasting message from client(%s) to all clients in room(%s)", clientID, reqClient.RoomID)
		for id, client := range clients {
			if reqClient.RoomID != client.RoomID {
				continue
			}
			if err := writeClientMessage(id, msg); err != nil {
				log.Printf("error writing client(%s) message: %s", clientID, err)
				continue
			}
		}
	}

	log.Print(logMessage)
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprintf(w, "%s\n", logMessage); err != nil {
		log.Printf("error writing client http response: %s", err)
	}
}

func listUsersHttpHandler(w http.ResponseWriter, r *http.Request) {
	// list all connected clients
	w.WriteHeader(http.StatusOK)
	for clientID := range clients {
		if _, err := fmt.Fprintf(w, "client(%s)\n", clientID); err != nil {
			log.Printf("error writing client http response: %s", err)
		}
	}
}

func listRoomsHttpHandler(w http.ResponseWriter, r *http.Request) {
	// list all unique rooms
	uniqueRooms := make(map[string]bool)
	for _, client := range clients {
		uniqueRooms[client.RoomID] = true
	}

	w.WriteHeader(http.StatusOK)
	for roomID := range uniqueRooms {
		if _, err := fmt.Fprintf(w, "room(%s)\n", roomID); err != nil {
			log.Printf("error writing client http response: %s", err)
		}
	}
}
