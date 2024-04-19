package main

import (
	"bufio"
	"fmt"
	"os"
	peer "peerproject/pair"
	"strings"
	"time"
)

func sleep() {
	time.Sleep(time.Second)
}

func readInput() []string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.Fields(text)
}

func handlePeer(MyPeer *peer.Peer, action string) {
	fmt.Print("\u001B[92mEnter the IP and port of peer you want to ", action, " :\u001B[39m")
	input := readInput()
	if len(input) == 2 {
		if input[0] == "localhost" {
			input[0] = "127.0.0.1"
		}
		if MyPeer.Comm[input[0]+":"+input[1]] == nil {
			fmt.Println("\u001B[92mNo such connection :(\u001B[39m")
			return
		}
		switch action {
		case "handle":
			peer.WriteReadConnection(MyPeer.Comm["127.0.0.1"+":"+input[1]], MyPeer)
		case "close":
			MyPeer.Close(input[0] + ":" + input[1])
		default:
			fmt.Println("\u001B[92mInvalid action.\u001B[39m")
		}
	} else {
		fmt.Println("\u001B[92mInvalid input.\u001B[39m")
	}
}

func inputProg() { // \u001B[92m \u001B[39m
	var MyPeer peer.Peer
	fmt.Println("\u001B[92mWelcome on Teemo2P !\u001B[32m")
	for {
		fmt.Print("\u001B[92mEnter a command :\u001B[39m")
		command := readInput()
		if len(command) <= 0 {
			continue
		} else {
			switch command[0] {
			case "launch a peer", "lp":
				fmt.Print("\u001B[92mGot it, Give me his IP & Port :\u001B[32m")
				input := readInput()
				if len(input) == 0 {
					MyPeer = peer.StartPeer("localhost", "3000", "online")
				} else if len(input) >= 2 {
					MyPeer = peer.StartPeer(input[0], input[1], "online")
				} else if len(input) == 1 {
					MyPeer = peer.StartPeer("localhost", input[0], "online")
				} else {
					fmt.Println("\u001B[92mMissing a field.\u001B[32m")
				}
			case "co", "connect":
				if !MyPeer.IsEmpty() {
					fmt.Print("\u001B[92mEnter the IP and port of peer you want to connect to:\u001B[32m")
					input := readInput()
					if len(input) == 2 {
						MyPeer.ConnectTo(input[0], input[1])
					} else if len(input) == 1 {
						MyPeer.ConnectTo("localhost", input[0])
					} else {
						fmt.Println("\u001B[92mInvalid input.\u001B[32m")
					}
				} else {
					fmt.Println("\u001B[92mYou need to launch a peer first.\u001B[32m")
				}
			case "hd", "handle":
				handlePeer(&MyPeer, "handle")
			case "cl", "close":
				handlePeer(&MyPeer, "close")
			case "exit":
				fmt.Println("\u001B[92mEnding the program :(, before we are closing all connections.\u001B[32m")
				for key, _ := range MyPeer.Comm {
					MyPeer.Close(key)
				}
				os.Exit(1)
			default:
				fmt.Println("\u001B[92mCommand not found here the list: lp (launch a peer), handle, close, exit\u001B[32m")
			}
		}
	}
}

func main() {
	inputProg()
	// P1 := peer.StartPeer("Peer1", "localhost", "3003", "online", make([]tools.File, 0))
	// sleep()
	// P2 := peer.StartPeer("Peer2", "localhost", "3007", "online", make([]tools.File, 0))
	// sleep()
	// go P2.ConnectTo(P1)
	select {}
}
