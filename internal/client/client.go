package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"wantsome.ro/messagingapp/internal/config"
	"wantsome.ro/messagingapp/pkg/models"
)

var currentRoom string

func RunClient() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run client.go [username]")
	}
	username := os.Args[1]
	config := config.LoadConfig()
	url := "ws://" + config.URL + ":" + config.Port + "/ws"

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("error dialing %s\n", err)
	}
	defer c.Close()

	// Send initial login message to store the users
	loginMessage := createLoginMessage(username)
	err = c.WriteJSON(loginMessage)
	if err != nil {
		log.Fatalf("error sending login message: %s\n", err)
	}

	done := make(chan bool)

	reader := bufio.NewReader(os.Stdin)

	// Handle incoming messages in a separate goroutine
	go handleServerMessages(c, username, done)

	// Show the main menu
	showMainMenu(reader, c, username)

	<-done
}

func showMainMenu(reader *bufio.Reader, c *websocket.Conn, username string) {
	for {
		fmt.Println("1. List users")
		fmt.Println("2. Select user to chat")
		fmt.Println("3. Create/Join a room")
		fmt.Println("4. Send message in current room")
		fmt.Println("5. Exit")

		fmt.Print("Enter choice: \n")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			// Implement the logic to list users
			requestUserList(c)
		case "2":
			// Implement the logic to chat with a selected user
			chatWithUser(reader, c, username)
		case "3":
			// Implement the logic to join a room
			joinRoom(reader, c, username)
		case "4":
			// Implement the logic to send a message in the current room
			sendMessageInRoom(reader, c, username)
		case "5":
			// Exit the program or return to a higher menu
			os.Exit(0)
		default:
			// Handle invalid choice
			fmt.Println("Invalid choice, please try again.")
		}
	}
}

func chatWithUser(reader *bufio.Reader, c *websocket.Conn, username string) {
	fmt.Print("Enter recipient's name: ")
	recipient, _ := reader.ReadString('\n')
	recipient = strings.TrimSpace(recipient)
	fmt.Print("Enter message (or type '/exit' to return to main menu): \n")
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		// Return to main menu if the user types /exit
		if text == "/exit" {
			break
		}

		message := createRoomMessage(username, text, "", recipient)

		err := c.WriteJSON(message)
		if err != nil {
			log.Printf("error writing %s\n", err)
			break
		}
	}
}

func handleServerMessages(c *websocket.Conn, username string, done chan bool) {
	// reading server messages
	go func() {
		defer close(done)
		for {
			_, messageBytes, err := c.ReadMessage()
			if err != nil {
				log.Printf("error reading: %s\n", err)
				return
			}

			var msg models.Message
			if err := json.Unmarshal(messageBytes, &msg); err != nil {
				log.Printf("error decoding message: %s\n", err)
				continue
			}

			switch msg.Type {
			case "user_list":
				fmt.Println("Active Users:", msg.Message)
			default:
				// Display the message if it's for the current room or a direct message to the user
				if msg.Room == currentRoom && msg.Recipient == "" && msg.UserName != username {
					fmt.Printf("[%s] %s: %s\n", msg.Room, msg.UserName, msg.Message)
				}
				if msg.Recipient == username && msg.UserName != username {
					fmt.Printf("[Private] %s: %s\n", msg.UserName, msg.Message)
				}
			}
		}
	}()
}

func requestUserList(c *websocket.Conn) {
	request := models.Message{
		Type: "list_users",
	}
	err := c.WriteJSON(request)
	if err != nil {
		log.Printf("error sending list_users request: %s\n", err)
	}
}

func joinRoom(reader *bufio.Reader, c *websocket.Conn, username string) {
	fmt.Print("Enter room number: ")
	room, _ := reader.ReadString('\n')
	room = strings.TrimSpace(room)

	currentRoom = room // Update the current room

	message := createJoinRoomMessage(username, room)

	c.WriteJSON(message)
}

func sendMessageInRoom(reader *bufio.Reader, c *websocket.Conn, username string) {
	if currentRoom == "" {
		fmt.Println("You are not in a room. Please join a room first.")
		return
	}
	fmt.Print("Enter message (or type '/exit' to return to main menu): \n")
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "/exit" {
			break
		}

		message := createRoomMessage(username, text, currentRoom, "")
		err := c.WriteJSON(message)
		if err != nil {
			log.Printf("error writing %s\n", err)
			break
		}

	}
}

// Function to create a login message
func createLoginMessage(username string) models.Message {
	return models.Message{
		Type:     "login",
		UserName: username,
	}
}

// Function to create a message for a room
func createRoomMessage(username, message, room, recipient string) models.Message {
	return models.Message{
		UserName:  username,
		Message:   message,
		Room:      room,
		Recipient: recipient,
	}
}

// Function to create a join_room message
func createJoinRoomMessage(username, room string) models.Message {
	return models.Message{
		Type:     "join_room",
		UserName: username,
		Room:     room,
	}
}
