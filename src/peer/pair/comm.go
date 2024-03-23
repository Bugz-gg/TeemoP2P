package peer_package

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func (p *Peer) HelloTrack(t Peer) {
	message := "annonce listen " + p.Port + " seed ["
	for _, valeur := range p.Files {
		name, size, pieceSize, key, isEmpty := valeur.GetFile()
		if isEmpty {
			message += fmt.Sprintf(`Name: %s, Size: %d, PieceSize: %d, Key: %s`, name, size, pieceSize, key)
		} else {
			break
		}
	}
	message += "]\n"
	conn, err := net.Dial("tcp", t.IP+":"+t.Port)
	errorCheck(err)
	// defer conn.Close()
	message = string(message)
	_, err = conn.Write([]byte(message))
	errorCheck(err)
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	errorCheck(err)
	fmt.Print("< ", string(buffer[:n]))
	p.Comm["tracker"] = conn
}

func WriteReadConnection(conn net.Conn) {
	// print(conn)
	// defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Waiting for input... :")
	message, _ := reader.ReadString('\n')
	fmt.Println(conn.LocalAddr().String(), ": < ", message)
	if message == "exit\n" {
		os.Exit(1)
	}
	conn.Write([]byte(message))

	buffer := make([]byte, 256)
	_, nerr := conn.Read(buffer)
	errorCheck(nerr)

}

func (p *Peer) ConnectTo(IP string, Port string) {
	conn, err := net.Dial("tcp", IP+":"+Port)
	errorCheck(err)
	// defer conn.Close()
	p.Comm[conn.RemoteAddr().String()] = conn
	fmt.Println(conn.LocalAddr(), " is connected to ", conn.RemoteAddr())
	WriteReadConnection(conn)
	// go handleConnection(conn)
	//Handle the response here !
	// fmt.Print(string(buffer[:n]))
}
