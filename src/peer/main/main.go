package main

import (
	peer "peerproject/pair"
	"peerproject/tools"
	"time"
)

func sleep() {
	time.Sleep(time.Second)
}
func main() {
	P1 := peer.StartPeer("Peer1", "localhost", "3003", "online", make([]tools.File, 0))
	sleep()
	P2 := peer.StartPeer("Peer2", "localhost", "3007", "online", make([]tools.File, 0))
	sleep()
	go P2.ConnectTo(P1)
	select {}
}
