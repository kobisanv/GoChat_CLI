// client.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// sendMessage sends user-input messages to the server.
func sendMessage(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)

	for {
		// Prompt the user to enter a message
		fmt.Print("Enter message: ")

		// Read the message from the console
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		// Send the message to the server
		conn.Write([]byte(message))
	}
}

func main() {
	// Prompt the user to enter their username
	fmt.Print("Enter your username: ")
	username, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatal("Error reading username:", err)
	}

	// Remove the newline character from the username
	username = username[:len(username)-1]

	// Connect to the server running on localhost:8080
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	// Send the username to the server
	conn.Write([]byte(username + "\n"))

	// Receive the chat room name from the server
	chatroomName, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal("Error reading chat room name:", err)
	}
	chatroomName = strings.TrimSpace(chatroomName)

	// Print a message indicating that the client is connected to the chat room
	fmt.Printf("Connected to the chat room: %s\n", chatroomName)

	// Start a goroutine to handle sending messages to the server
	go sendMessage(conn)

	// Create a reader to receive messages from the server
	reader := bufio.NewReader(conn)
	for {
		// Read a message from the server
		message, err := reader.ReadString('\n')
		if err != nil {
			// If there's an error, assume the server closed the connection
			log.Println("Server closed connection")
			return
		}

		// Print the received message to the console
		fmt.Print(message)
	}
}
