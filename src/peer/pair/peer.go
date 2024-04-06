package peer_package

import (
	"fmt"
	"net"
	"time"

	"gopkg.in/ini.v1"

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
	file, err := ini.Load("config.ini")
	errorCheck(err)
	var track Peer

	section := file.Section("Tracker")
	track.IP = section.Key("ip").String()
	track.Port = section.Key("port").String()
	track.Status = "online"
	track.Type = "tracker"
	return track
}

func StartPeer(IP string, Port string, Type string) Peer {
	track := GetConfig()
	peer := Peer{
		IP:    IP,
		Port:  Port,
		Type:  Type,
		Files: tools.FillStructFromConfig(),
		Comm:  make(map[string]net.Conn),
	}
	go peer.startListening()
	time.Sleep(time.Second)
	peer.HelloTrack(track)
	return peer
}
