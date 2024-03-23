package peer_package

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	tools "peerproject/tools"
)

type Peer struct {
	IP     string
	Port   string
	Status string
	Type   string
	Files  []tools.File
	Comm   map[string]net.Conn
}

func (p *Peer) IsEmpty() bool {
	if len(p.IP) == 0 || len(p.Port) == 0 {
		return true
	} else {
		return false
	}
}
func errorCheck(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
}

func GetConfig() Peer {
	file, err := os.Open("./config.ini")
	errorCheck(err)
	var track Peer

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
				track.IP = value
			case "tracker-port":
				track.Port = value
			}
		}
	}
	track.Status = "online"
	track.Type = "tracker"
	return track
}

func StartPeer(IP string, Port string, Type string, Files []tools.File) Peer {
	track := GetConfig()
	peer := Peer{
		IP:    IP,
		Port:  Port,
		Type:  Type,
		Files: Files,
		Comm:  make(map[string]net.Conn),
	}
	go peer.startListening()
	time.Sleep(time.Second)
	peer.HelloTrack(track)
	return peer
}
