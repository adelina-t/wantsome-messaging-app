
# Documentation for Messaging App

## Overview
This project comprises several Go files that work together to create a messaging application with a client-server architecture. The application uses WebSockets for real-time communication.

## Files Description

### handler.go
- **Package**: server
- **Purpose**: Manages WebSocket connections for the messaging application.
- **Functionality**:
  - Handles HTTP requests and upgrades them to WebSocket connections.
  - Manages user connections and chat rooms.
  - Broadcasts messages to appropriate recipients.

### server.go
- **Package**: server
- **Purpose**: Runs the server side of the messaging application.
- **Functionality**:
  - Loads server configuration from environment variables.
  - Registers HTTP handlers for the root path and WebSocket connections.
  - Starts the HTTP server and listens for shutdown signals.

### client.go
- **Package**: client
- **Purpose**: Manages the client-side operations of the messaging application.
- **Functionality**:
  - Connects to the server via WebSocket.
  - Handles sending and receiving of messages.
  - Supports command-line arguments for user customization.

### config.go
- **Package**: config
- **Purpose**: Provides configuration settings for the server.
- **Functionality**:
  - Defines `Config` struct with URL and Port fields.
  - Loads configuration from environment variables with default fallbacks.

### models.go
- **Package**: models
- **Purpose**: Defines data structures used in the application.
- **Functionality**:
  - Includes `Message` struct for representing chat messages.

## Integration Tests
Integration tests will be written to test the interaction between the client and server components, focusing on:
- Connection handling.
- Message sending and receiving.
- Error handling and edge cases.
