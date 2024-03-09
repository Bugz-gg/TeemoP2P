package peer

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Peer struct {
	ID      string
	Address string
}

func errorCheck(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Connection established from", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		fmt.Print("Received message:", message)

		if strings.TrimSpace(message) == "exit" {
			return
		}

		fmt.Print(">")
		response, _ := reader.ReadString('\n')
		conn.Write([]byte(response))
	}
}

func (p *Peer) Start(targetAddr string, clientMode bool) {
	if clientMode {
		// Peer acts as a client, initiates connection to target address
		targetTCPAddr, err := net.ResolveTCPAddr("tcp", targetAddr)
		errorCheck(err)

		conn, err := net.DialTCP("tcp", nil, targetTCPAddr)
		errorCheck(err)

		fmt.Printf("Peer %s connected to %s\n", p.ID, targetAddr)

		go handleConnection(conn)
	} else {
		// Peer acts as a server, listens for incoming connections
		listening, err := net.Listen("tcp", p.Address)
		errorCheck(err)
		defer listening.Close()

		fmt.Printf("Peer %s listening on %s\n", p.ID, p.Address)

		for {
			conn, err := listening.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}
			go handleConnection(conn)
		}
	}
}
