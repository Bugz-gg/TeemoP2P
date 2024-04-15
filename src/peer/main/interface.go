package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

func main() {
	http.HandleFunc("/", indexHandler)
	// http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/launch", launchHandler)
	// http.HandleFunc("/connect", connectHandler)
	// http.HandleFunc("/handle", handleHandler)
	// http.HandleFunc("/close", closeHandler)
	// http.HandleFunc("/exit", exitHandler)

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("./main/index.html"))
	tpl.Execute(w, nil)
}

func launchHandler(w http.ResponseWriter, r *http.Request) {
	// Generate a random port number
	port := strconv.Itoa(3000 + rand.Intn(1000)) // Generate a random port in the range 3000-3999

	// Start the peer with the generated port
	// peer.StartPeer("localhost", port, "online")
	fmt.Println("Peer launching... Port:", port)
	// time.Sleep(2 * time.Second)
	// Simulate a successful response
	// successMessage := fmt.Sprintf("Peer launched successfully! IP: %s, Port: %s", MyPeer.IP, MyPeer.Port)

	// Send the success message as the response
	fmt.Fprint(w, port)
	// w.WriteHeader(http.StatusOK)
}

func connectHandler(w http.ResponseWriter, r *http.Request) {
	// Implement connecting to a peer
	// Example:
	// fmt.Fprint(w, "Connected to peer successfully!")
}

func handleHandler(w http.ResponseWriter, r *http.Request) {
	// Implement handling a peer
	// Example:
	// fmt.Fprint(w, "Peer handled successfully!")
}

func closeHandler(w http.ResponseWriter, r *http.Request) {
	// Implement closing a connection
	// Example:
	// fmt.Fprint(w, "Connection closed successfully!")
}

func exitHandler(w http.ResponseWriter, r *http.Request) {
	// Implement exiting the program
	// Example:
	// fmt.Fprint(w, "Exiting the program...")
}
