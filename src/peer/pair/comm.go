package peer_package

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func (p *peer) HelloTrack(t peer) {
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
	defer conn.Close()
	message = string(message)
	_, err = conn.Write([]byte(message))
	errorCheck(err)
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	errorCheck(err)
	fmt.Print("< ", string(buffer[:n]))
}

func (p *peer) ConnectTo(s peer) {
	conn, err := net.Dial("tcp", s.IP+":"+s.Port)
	errorCheck(err)
	defer conn.Close()
	fmt.Println(p.Name, " is connected to ", s.Name)
	user := p.Name
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Waiting for input... :")
		message, _ := reader.ReadString('\n')
		fmt.Println(user, ": < ", message)
		if message == "exit\n" {
			os.Exit(1)
		}
		conn.Write([]byte(message))

		buffer := make([]byte, 256)
		n, err := conn.Read(buffer)
		errorCheck(err)
		fmt.Print(string(buffer[:n]))
	}
}
