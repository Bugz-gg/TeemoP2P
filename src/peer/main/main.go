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
	fmt.Print("Enter the IP and port of peer you want to ", action, " :")
	input := readInput()
	if len(input) == 2 {
		if input[0] == "localhost" {
			input[0] = "127.0.0.1"
		}
		if MyPeer.Comm[input[0]+":"+input[1]] == nil {
			fmt.Println("No such connection :(")
			return
		}
		switch action {
		case "handle":
			peer.WriteReadConnection(MyPeer.Comm["127.0.0.1"+":"+input[1]], MyPeer)
		case "close":
			MyPeer.Close(input[0] + ":" + input[1])
		default:
			fmt.Println("Invalid action.")
		}
	} else {
		fmt.Println("Invalid input.")
	}
}

func inputProg() {
	var MyPeer peer.Peer
	fmt.Println("Welcome on Teemo2P !")
	for {
		fmt.Print("Enter a command : ")
		command := readInput()
		if len(command) <= 0 {
			continue
		} else {
			switch command[0] {
			case "launch a peer", "lp":
				fmt.Print("Got it, Give me his IP & Port : ")
				input := readInput()
				if len(input) == 0 {
					MyPeer = peer.StartPeer("localhost", "3000", "online")
				} else if len(input) >= 2 {
					MyPeer = peer.StartPeer(input[0], input[1], "online")
				} else if len(input) == 1 {
					MyPeer = peer.StartPeer("localhost", input[0], "online")
				} else {
					fmt.Println("Missing a field.")
				}
			case "co", "connect":
				if !MyPeer.IsEmpty() {
					fmt.Print("Enter the IP and port of peer you want to connect to: ")
					input := readInput()
					if len(input) == 2 {
						MyPeer.ConnectTo(input[0], input[1])
					} else if len(input) == 1 {
						MyPeer.ConnectTo("localhost", input[0])
					} else {
						fmt.Println("Invalid input.")
					}
				} else {
					fmt.Println("You need to launch a peer first.")
				}
			case "hd", "handle":
				handlePeer(&MyPeer, "handle")
			case "cl", "close":
				handlePeer(&MyPeer, "close")
			case "exit":
				fmt.Println("Ending the program :(")
				os.Exit(1)
			default:
				fmt.Println("Command not found here the list: lp (launch a peer), handle, close, exit")
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
