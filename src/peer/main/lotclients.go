package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"time"
)

func main() {
	// Number of clients to spawn
	numClients := 100000

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Loop to spawn multiple clients
	for i := 0; i < numClients; i++ {
		go func(clientNum int) {
			// Connect to the server
			conn, err := net.Dial("tcp", "localhost:3000")
			if err != nil {
				fmt.Printf("Client %d: Error connecting: %v\n", clientNum, err)
				return
			}
			defer conn.Close()

			// Start a separate goroutine to listen for messages from the server
			go func() {
				for {
					message, err := bufio.NewReader(conn).ReadString('\n')
					if err != nil {
						fmt.Printf("Client %d: Error reading: %v\n", clientNum, err)
						return
					}
					fmt.Printf("Client %d received: %s", clientNum, message)
				}
			}()

			// Generate and send random messages to the server
			for {
				message := generateRandomMessage()

				// Send the message to the server
				_, err := fmt.Fprintf(conn, message+"\n")
				if err != nil {
					fmt.Printf("Client %d: Error sending: %v\n", clientNum, err)
					return
				}

				// Sleep for a random duration before sending the next message
			}
		}(i + 1) // Pass client number to the anonymous function
	}

	// Keep the main goroutine alive
	select {}
}

// Function to generate a random message
func generateRandomMessage() string {
	messages := []string{
		"Hello",
		"Hi there",
		"Random message",
		"Greetings",
		"How are you?",
		"Goodbye",
		"See you later",
		"Nice to meet you",
	}
	return messages[rand.Intn(len(messages))]
}
