package peer_package

import (
	"fmt"
	"net"
	"time"

	tools "peerproject/tools"
)

type Peer struct {
	IP     string
	Port   string
	Status string
	Type   string
	Files  map[string]*tools.File
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
		// panic(err)
	}
}

func GetConfig() Peer {
	/*file, err := ini.Load("config.ini")
	errorCheck(err)
	if err != nil {
		fmt.Printf("\033[0;31mYou have no config.ini !\033[0m\n")
	}*/
	var track Peer

	//section := file.Section("Tracker")
	//track.IP = section.Key("ip").String()
	//track.Port = section.Key("port").String()
	track.IP = tools.GetValueFromConfig("Tracker", "ip")
	track.Port = tools.GetValueFromConfig("Tracker", "port")
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
		Files: make(map[string]*tools.File),
		Comm:  make(map[string]net.Conn),
	}
	peer.Files = tools.FillFilesFromConfig(IP + ":" + Port)
	/*if peer.Files == nil {
		peer.Files = make(map[string]*tools.File)
	}*/
	tools.LocalFiles = &peer.Files
	go peer.startListening()
	go peer.sendupdate()
	time.Sleep(time.Millisecond)
	peer.HelloTrack(track)
	//go peer.rarepiece()
	return peer
}
