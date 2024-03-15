package main

import (
	peer "peerproject/pair"
	"peerproject/tools"
	"time"
)

func main() {
	P1 := peer.StartPeer("Peer1", "localhost", "3003", "online", make([]tools.File, 0))
	time.Sleep(time.Second)
	P2 := peer.StartPeer("Peer2", "localhost", "3007", "online", make([]tools.File, 0))
	go P2.ConnectTo(P1)
	select {}
}
