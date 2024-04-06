package peer_package

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"peerproject/tools"
	"strings"
	// "peerproject/tools"
)

func (p *Peer) HelloTrack(t Peer) {
	message := "announce listen " + p.Port + " seed ["
	for _, valeur := range p.Files {
		name, size, pieceSize, key, isEmpty := valeur.GetFile()
		if isEmpty {
			message += fmt.Sprintf(`%s %d %d %s `, name, size, pieceSize, key)
		} else {
			break
		}
	}
	message = strings.TrimSuffix(message, " ")
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
	fd, nerr := conn.Read(buffer)
	errorCheck(nerr)
	if fd > 0 {
		mess := string(buffer)
		switch mess {
		case "data", "data\n":
			_, _ = tools.DataCheck(mess)
		case "have", "have\n":
			_, _ = tools.HaveCheck(mess)
		case "ok", "ok\n":
			fmt.Println("> ", mess)
		case "list", "list\n":
			_, _ = tools.ListCheck(mess)
		case "peers", "peers\n":
			_, _ = tools.PeersCheck(mess)
		default:
			panic("valeur par default et pas parmi la liste")

		}
	}
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
