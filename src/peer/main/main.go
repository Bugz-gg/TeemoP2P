package main

import (
	"bufio"
	"fmt"
	"os"
	peer "peerproject/pair"
)

func main() {
	// Get the target address from user input for Peer1
	fmt.Print("Enter the target IP and port for Peer1 (e.g., localhost:3001): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	targetAddr1 := scanner.Text()

	// Get the target address from user input for Peer2
	fmt.Print("Enter the target IP and port for Peer2 (e.g., localhost:3000): ")
	scanner.Scan()
	targetAddr2 := scanner.Text()

	peer1 := peer.Peer{ID: "Peer1", Address: "localhost:3000"}
	peer2 := peer.Peer{ID: "Peer2", Address: "localhost:3001"}

	// Peer1 initiates the connection to Peer2
	go peer1.Start(targetAddr2, true) // Pass true to indicate peer1 is a client

	// Peer2 listens for incoming connections
	go peer2.Start(targetAddr1, false) // Pass false to indicate peer2 is a server

	select {}
}
