package peer_package

import (
	"fmt"
	"io"
	"net"
	"peerproject/tools"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

func (p *Peer) Close(t string) {
	fmt.Println("Carefull the connection with", t, "is now closed")
	p.Comm[t].Close()
	delete(p.Comm, t)
}

func ReadResConnection(conn net.Conn, sem chan struct{}) {
	defer conn.Close()
	defer func() { <-sem }() // Release semaphore after finishing

	inidata, err := ini.Load("./config.ini")
	errorCheck(err)

	section := inidata.Section("Peer")
	max_attempts, err := strconv.Atoi(section.Key("max_message_attempts").String())
	errorCheck(err)

	for {
		if max_attempts == 0 {
			fmt.Println("To many tries connection closed")
			return
		}
		buffer := make([]byte, 256)
		n, err := conn.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Reading error:", err)
			return
		}
		if n > 0 {
			mess := string(buffer[:n])
			word := strings.Fields(mess)
			switch word[0] {
			case "interested", "interested\n":
				_, _ = tools.InterestedCheck(mess)
			case "getpieces", "getpieces\n":
				_, _ = tools.GetPiecesCheck(mess)
			case "have", "have\n":
				_, _ = tools.HaveCheck(mess)
			case "exit", "exit\n":
				return
			default:
				max_attempts--
				// fmt.Print(conn.LocalAddr(), "> ", string(buffer[:n]))
				// _, err := conn.Write(buffer[:n])
				// errorCheck(err)
				fmt.Println("You have ", max_attempts, " tries remaining")
			}
		}
	}
}

func (p *Peer) startListening() {
	l, err := net.Listen("tcp", p.IP+":"+p.Port)
	errorCheck(err)
	defer l.Close()

	fmt.Println("Server listening on port", p.Port)
	inidata, err := ini.Load("./config.ini")
	errorCheck(err)

	section := inidata.Section("Peer")
	maxConcurrent, err := strconv.Atoi(section.Key("max_concurrence").String())
	errorCheck(err)
	sem := make(chan struct{}, maxConcurrent)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Acceptation error:", err)
			continue
		}

		// Acquire semaphore
		sem <- struct{}{}

		// Start goroutine with limited concurrency
		go ReadResConnection(conn, sem)
	}
}
