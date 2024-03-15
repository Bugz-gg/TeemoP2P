package peer_package

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	tools "peerproject/tools"
)

type peer struct {
	Name   string
	IP     string
	Port   string
	Status string
	Type   string
	Files  []tools.File
}

func errorCheck(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
}

func GetConfig() peer {
	file, err := os.Open("./config.ini")
	errorCheck(err)
	var track peer

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

func StartPeer(Name string, IP string, Port string, Type string, Files []tools.File) peer {
	track := GetConfig()
	peer := peer{
		Name:  Name,
		IP:    IP,
		Port:  Port,
		Type:  Type,
		Files: Files,
	}
	go peer.startListening()
	time.Sleep(time.Second)
	go peer.HelloTrack(track)
	return peer
}
