package peer_package

import (
	"fmt"
	"io"
	"net"
)

func (p *Peer) Close(t string) {
	fmt.Println("Carefull the connection with", t)
	p.Comm[t].Close()
}

func ReadResConnection(conn net.Conn) {
	// defer conn.Close()

	for {
		buffer := make([]byte, 256)
		n, err := conn.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Reading error:", err)
			return
		}
		if n > 0 {
			fmt.Print(conn.LocalAddr(), "> ", string(buffer[:n]))
			_, err := conn.Write(buffer[:n])
			errorCheck(err)
			// if err != nil {
			// 	fmt.Println("Writing error:", err)
			// 	return
			// }
		}
	}
}

func (p *Peer) startListening() {

	l, err := net.Listen("tcp", p.IP+":"+p.Port)
	errorCheck(err)
	defer l.Close()

	fmt.Println("Server listening on port", p.Port)

	for {
		conn, err := l.Accept()
		p.Comm[conn.RemoteAddr().String()] = conn
		if err != nil {
			fmt.Println("Acceptation error:", err)
			continue
		}

		// Démarrez une goroutine pour gérer la connexion
		go ReadResConnection(conn)
	}
}
