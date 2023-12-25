// server.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// broadcastMessage sends a message to all clients except the sender.
// This function is used to broadcast messages to all connected clients.
func broadcastMessage(message string, clients map[net.Conn]string, sender net.Conn) {
	for client := range clients {
		// Exclude the sender from the broadcast
		if client != sender {
			client.Write([]byte(message + "\n"))
		}
	}
}

// handleConnection handles the communication with a connected client.
func handleConnection(conn net.Conn, clients map[net.Conn]string, chatroomName string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	// Read the username from the client
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading username:", err)
		return
	}

	// Trim leading and trailing whitespaces from the username
	username = strings.TrimSpace(username)
	fmt.Printf("[%s] joined the chat room\n", username)

	// Broadcast a message to all clients about the new user joining
	broadcastMessage(fmt.Sprintf("[%s] joined the chat room", username), clients, conn)

	// Store the new client in the clients map
	clients[conn] = username

	for {
		// Read the message from the client
		message, err := reader.ReadString('\n')
		if err != nil {
			// If there's an error, assume the user left and broadcast a message
			log.Printf("[%s] left the chat room\n", username)
			delete(clients, conn)
			broadcastMessage(fmt.Sprintf("[%s] left the chat room", username), clients, conn)
			return
		}

		// Trim leading and trailing whitespaces from the message
		message = strings.TrimSpace(message)
		// Print the message to the server console
		fmt.Printf("{%s} sent a message to {%s}: %s\n", username, chatroomName, message)

		// Broadcast the message to all clients
		broadcastMessage(fmt.Sprintf("{%s} sent a message to {%s}: %s", username, chatroomName, message), clients, conn)
	}
}

func main() {
	// Prompt the user to enter the chat room name
	fmt.Print("Enter chat room name: ")
	chatroomName, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatal("Error reading chat room name:", err)
	}

	// Trim leading and trailing whitespaces from the chat room name
	chatroomName = strings.TrimSpace(chatroomName)

	// Create a listener for incoming connections on port 8080
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error creating server:", err)
	}
	defer ln.Close()

	// Print a message indicating that the server is running with the chat room name
	fmt.Printf("Server is up and running. Chat room: '%s'\n", chatroomName)

	// Create a map to store connected clients
	clients := make(map[net.Conn]string)

	// Continuously accept and handle incoming connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go handleConnection(conn, clients, chatroomName)
	}
}
