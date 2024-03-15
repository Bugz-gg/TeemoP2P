package peer_package

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buffer := make([]byte, 256)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Reading error:", err)
			return
		}

		if n > 0 {
			fmt.Println(conn.LocalAddr(), " > ", string(buffer[:n]))
			_, err := conn.Write(buffer[:n])
			errorCheck(err)
		}
	}
}

func (p *peer) startListening() {

	l, err := net.Listen("tcp", p.IP+":"+p.Port)
	errorCheck(err)
	defer l.Close()

	fmt.Println("Server listening on port", p.Port)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Acceptation error:", err)
			continue
		}

		// Démarrez une goroutine pour gérer la connexion
		go handleConnection(conn)
	}
}
