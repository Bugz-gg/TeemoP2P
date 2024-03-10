package main

import (
	peer "peerproject/pair"
)

func main() {
	track := peer.GetConfig()
	go peer.StartListening(track.IP, track.Port)
	peer1 := peer.Peer{IP: "localhost", Port: "3005", Status: "peer"}
	go peer1.StartPeer()
	select {}
}
