package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"wantsome.ro/messagingapp/pkg/models"
)

func TestLogRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	loggedHandler := logRequest(handler)

	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	loggedHandler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "OK", rr.Body.String())
}

func TestHomeHttpHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	homeHttpHandler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "Welcome to wantsome.ro messaging app!\n", rr.Body.String())
}

func TestGetClientIDHeaderMissing(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	_, err = getClientIDHeader(nil, req)
	require.Equal(t, "Client-ID header is missing", err.Error())
}

func TestGetRoomIDHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	req.Header.Set("Room-ID", "default")
	roomID, err := getRoomIDHeader(nil, req)
	require.NoError(t, err)
	require.Equal(t, "default", roomID)
}

func TestGetRoomIDHeaderMissing(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	_, err = getRoomIDHeader(nil, req)
	require.Equal(t, "Room-ID header is missing", err.Error())
}

func TestClientWsHandler(t *testing.T) {
	// create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(clientWsHandler))
	defer server.Close()

	// convert http://127.0.0.1 to ws://127.0.0.1
	u := "ws" + strings.TrimPrefix(server.URL, "http")

	// create a header with the client and room IDs
	header := http.Header{}
	header.Add("Client-ID", "testClient")
	header.Add("Room-ID", "default")

	// connect to the server with the header
	ws, _, err := websocket.DefaultDialer.Dial(u, header)
	require.NoError(t, err)
	defer ws.Close()

	// write message to server
	err = ws.WriteMessage(websocket.TextMessage, []byte("Hello everyone!"))
	require.NoError(t, err)

	// read response and unmarshal into a Message
	_, p, err := ws.ReadMessage()
	require.NoError(t, err)

	var msg models.Message
	err = json.Unmarshal(p, &msg)
	require.NoError(t, err)

	// check the message
	expected := models.Message{
		Content: fmt.Sprintf("Hello %s! Welcome to chat room %s!", "testClient", "default"),
	}
	require.Equal(t, expected, msg)
}

func TestSendMessageHttpHandlerBadRequest(t *testing.T) {
	handler := http.HandlerFunc(sendMessageHttpHandler)

	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequest("POST", server.URL, strings.NewReader(`{"content":"test"}`))
	require.NoError(t, err)

	req.Header.Set("Client-ID", "default")
	req.Header.Set("Room-ID", "default")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Equal(t, "error: client(default) not connected\n", rr.Body.String())
}

func TestListUsersHttpHandler(t *testing.T) {
	handler := http.HandlerFunc(listUsersHttpHandler)

	server := httptest.NewServer(handler)
	defer server.Close()

	// Add a user
	clients["test"] = &models.Client{
		RoomID: "default",
	}

	req, err := http.NewRequest("GET", server.URL, nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "client(test)\n", rr.Body.String())
}
