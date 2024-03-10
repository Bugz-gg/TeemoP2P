package peer

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Peer struct {
	IP     string
	Port   string
	Status string
}

func errorCheck(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func GetConfig() Peer {
	file, err := os.Open("./config.ini")
	errorCheck(err)
	var peer Peer

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)

		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch key {
			case "tracker-address":
				peer.IP = value
			case "tracker-port":
				peer.Port = value
			}
		}
	}
	peer.Status = "tracker"
	return peer
}

func (p *Peer) StartListening() {

	// peers := make(map[string]Peer)
	server, err := net.Listen("tcp", p.IP+":"+p.Port)
	errorCheck(err)

	defer server.Close()
	fmt.Println("Listening on " + p.IP + ":" + p.Port)
	fmt.Println("Waiting for client...")

	conn, err := server.Accept()
	errorCheck(err)
	defer conn.Close()
	if p.Status == "tracker" {
		//tools de kevin
	}
	for {
		buffer := make([]byte, 256)
		n, err := conn.Read(buffer)
		errorCheck(err)

		if n > 0 {
			fmt.Print("<Message> ", string(buffer[:n]))
			_, err := conn.Write(buffer[:n])
			errorCheck(err)
		}
	}
}

func (p *Peer) StartPeer() {
	go p.StartListening()
	message := "< annonce listen " + p.Port + "[]"
	tracker := GetConfig()
	serv_tcp_addr, err := net.ResolveTCPAddr("tcp", tracker.IP+":"+tracker.Port)
	errorCheck(err)

	sockfd, err := net.DialTCP("tcp", nil, serv_tcp_addr)
	errorCheck(err)
	sockfd.Write([]byte(message))

}
