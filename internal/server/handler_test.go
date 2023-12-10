package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// 	if receivedMessage.Message != sendMessage.Message || receivedMessage.UserName != sendMessage.UserName {
// 		t.Errorf("Mismatch in sent and received messages. Expected: %v, Got: %v", sendMessage, receivedMessage)
// 	}
// }

func TestHome(t *testing.T) {
	// Set up a test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		home(w, r)
	}))
	defer testServer.Close()

	// Make a request to the home endpoint
	resp, err := http.Get(testServer.URL)
	if err != nil {
		t.Fatalf("Error making GET request: %s", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
