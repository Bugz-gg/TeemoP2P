package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"

	peer "peerproject/pair"
	tools "peerproject/tools"
)

var MyPeer peer.Peer

func sleep() {
	time.Sleep(time.Second)
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func handlePeer(MyPeer *peer.Peer, action string) {
	fmt.Printf("\u001B[92mEnter the number of the peer you want to %s:\u001B[39m\n", action)
	index := 0
	peerList := make([]string, 0, len(MyPeer.Comm))
	for key := range MyPeer.Comm {
		peerList = append(peerList, key)
		fmt.Printf("\u001B[96m(%d) %s\u001B[39m\n", index, key)
		index++
	}
	fmt.Print("\u001B[92mSelection:\u001B[39m ")
	input := readInput()
	peerNum, err := strconv.Atoi(input)
	if err != nil || peerNum < 0 || peerNum >= len(peerList) {
		fmt.Println("\u001B[91mInvalid selection.\u001B[39m")
		return
	}

	selectedPeer := peerList[peerNum]
	fmt.Println("\u001B[92mAvailable commands:\u001B[39m")
	if selectedPeer == peer.GetConfig().IP+":"+peer.GetConfig().Port || selectedPeer == "tracker" {
		fmt.Println("\u001B[96m(0) look\u001B[39m")
		fmt.Println("\u001B[96m(1) getfile\u001B[39m")
	} else {
		fmt.Println("\u001B[96m(2) interested\u001B[39m")
		fmt.Println("\u001B[96m(3) getpieces\u001B[39m")
	}

	fmt.Print("\u001B[92mEnter command number (press enter for default action):\u001B[39m ")
	commandInput := readInput()
	if commandInput == "" {
		switch action {
		case "handle":
			peer.WriteReadConnection(MyPeer.Comm[selectedPeer], MyPeer)
		case "close":
			MyPeer.Close(selectedPeer)
		default:
			fmt.Println("\u001B[91mInvalid action.\u001B[39m")
		}
		return
	}

	commandNum, err := strconv.Atoi(commandInput)
	if err != nil || commandNum < 0 || commandNum > 3 {
		fmt.Println("\u001B[91mInvalid command.\u001B[39m")
		return
	}

	switch commandNum {
	case 0:
		fmt.Println("\u001B[92mPerforming 'look' command...\u001B[39m")
		fmt.Print("Which criteria you are looking for: ")
		commandInput := readInput()
		if commandInput == "" {
			peer.WriteReadConnection(MyPeer.Comm[selectedPeer], MyPeer, "look []\n")
		} else {
			peer.WriteReadConnection(MyPeer.Comm[selectedPeer], MyPeer, "look ["+commandInput+"]\n")
		}
	case 1:
		fmt.Println("\u001B[92mPerforming 'getfile' command...\u001B[39m")
		remoteFileKeys := make([]string, 0, len(tools.RemoteFiles))
		for key := range tools.RemoteFiles {
			remoteFileKeys = append(remoteFileKeys, key)
		}

		for i, key := range remoteFileKeys {
			file := tools.RemoteFiles[key]
			fmt.Printf("(%d) %s %s\n", i, file.Name, key)
		}

		fmt.Print("\u001B[92mSelection:\u001B[39m ")
		input := readInput()
		fileNum, err := strconv.Atoi(input)
		if err != nil || fileNum < 0 || fileNum >= len(remoteFileKeys) {
			fmt.Println("\u001B[91mInvalid selection.\u001B[39m")
			return
		}

		selectedFile := remoteFileKeys[fileNum]

		fmt.Printf("\u001B[92mPerforming 'getfile' command for file %s...\u001B[39m\n", selectedFile)
		peer.WriteReadConnection(MyPeer.Comm[selectedPeer], MyPeer, fmt.Sprintf("getfile %s\n", selectedFile))

	case 2:
		fmt.Println("\u001B[92mPerforming 'interested' command...\u001B[39m")
		remoteFileKeys := make([]string, 0, len(tools.RemoteFiles))
		for key := range tools.RemoteFiles {
			remoteFileKeys = append(remoteFileKeys, key)
		}

		for i, key := range remoteFileKeys {
			file := tools.RemoteFiles[key]
			fmt.Printf("(%d) %s %s\n", i, file.Name, key)
		}

		fmt.Print("\u001B[92mSelection:\u001B[39m ")
		input := readInput()
		fileNum, err := strconv.Atoi(input)
		if err != nil || fileNum < 0 || fileNum >= len(remoteFileKeys) {
			fmt.Println("\u001B[91mInvalid selection.\u001B[39m")
			return
		}

		selectedFile := remoteFileKeys[fileNum]

		fmt.Printf("\u001B[92mPerforming 'interested' command for file %s...\u001B[39m\n", selectedFile)
		peer.WriteReadConnection(MyPeer.Comm[selectedPeer], MyPeer, fmt.Sprintf("interested %s\n", selectedFile))

	case 3:
		fmt.Println("\u001B[92mPerforming 'getpieces' command...\u001B[39m")
		remoteFileKeys := make([]string, 0, len(tools.RemoteFiles))
		for key := range tools.RemoteFiles {
			if _, valid := MyPeer.Files[key]; valid {
				remoteFileKeys = append(remoteFileKeys, key)
			}
		}

		// Display files with their keys
		for i, key := range remoteFileKeys {
			file := tools.RemoteFiles[key]
			fmt.Println(MyPeer.Comm[selectedPeer].RemoteAddr().String())
			fmt.Printf("(%d) %s %s (Total Pieces: %d)\n", i, file.Name, key, tools.BufferBitSize(*file))
		}

		fmt.Print("\u001B[92mSelection:\u001B[39m ")
		input := readInput()
		fileNum, err := strconv.Atoi(input)
		if err != nil || fileNum < 0 || fileNum >= len(remoteFileKeys) {
			fmt.Println("\u001B[91mInvalid selection.\u001B[39m")
			return
		}

		selectedFile := remoteFileKeys[fileNum]

		fmt.Printf("\u001B[92mPerforming 'getpieces' command for file %s...\u001B[39m\n", selectedFile)

		fmt.Print("Enter specific pieces (e.g., '3 5 7 8 9') or a range (e.g., '3-9'): ")
		piecesInput := readInput()
		var pieces []int
		if strings.Contains(piecesInput, "-") {
			rangeParts := strings.Split(piecesInput, "-")
			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				fmt.Println("\u001B[91mInvalid range.\u001B[39m")
				return
			}
			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				fmt.Println("\u001B[91mInvalid range.\u001B[39m")
				return
			}
			for i := start; i <= end; i++ {
				pieces = append(pieces, i)
			}
		} else {
			piecesStr := strings.Fields(piecesInput)
			for _, pieceStr := range piecesStr {
				piece, err := strconv.Atoi(pieceStr)
				if err != nil {
					fmt.Println("\u001B[91mInvalid input.\u001B[39m")
					return
				}
				pieces = append(pieces, piece)
			}
		}

		var currSize uint64 = 0
		var tmpIndexes []string
		buffSize, _ := strconv.ParseUint(tools.GetValueFromConfig("Peer", "max_buff_size"), 10, 64)

		for index := range pieces {
			tmpIndexes = append(tmpIndexes, strconv.Itoa(index))
			currSize += MyPeer.Files[selectedFile].PieceSize
			if currSize+100 >= buffSize {
				peer.ResponsesRemaining.Add(1)
				peer.WriteReadConnection(MyPeer.Comm[selectedPeer], MyPeer, "getpieces "+selectedFile+" ["+strings.Join(tmpIndexes, " ")+"]\n")
				tmpIndexes = []string{}
				currSize = 0
			}
		}
		if len(tmpIndexes) != 0 {
			peer.ResponsesRemaining.Add(1)
			peer.WriteReadConnection(MyPeer.Comm[selectedPeer], MyPeer, "getpieces "+selectedFile+" ["+strings.Join(tmpIndexes, " ")+"]\n")
		}
		go func() { peer.ChannSignal(&peer.ResponsesRemainingUpdated) }()
		peer.WaitFor(peer.DlDone, true)

		// command := fmt.Sprintf("getpieces %s %v\n", selectedFile, pieces)
		// peer.WriteReadConnection(MyPeer.Comm[selectedPeer], MyPeer, command)

	default:
		fmt.Println("\u001B[91mInvalid command.\u001B[39m")
	}
}

func inputProg() {
	//var MyPeer peer.Peer
	fmt.Println("\u001B[92mWelcome to Teemo2P!\u001B[39m")
	fmt.Print("\u001B[92mEnter IP & Port (e.g., localhost 3000):\u001B[39m ")
	input := readInput()
	fields := strings.Fields(input)
	switch len(fields) {
	case 0:
		MyPeer = peer.StartPeer("localhost", "3000", "online")
	case 1:
		MyPeer = peer.StartPeer("localhost", fields[0], "online")
	case 2:
		MyPeer = peer.StartPeer(fields[0], fields[1], "online")
	default:
		fmt.Println("\u001B[91mInvalid input.\u001B[39m")
	}

	availableCommands := []int{0, 4, 5, 6} // Handle, Connect, Close, Exit
	for {
		peer.ResponsesRemaining.Store(0)
		var downloadables []tools.File
		for _, file := range tools.RemoteFiles {
			if len(file.Peers) > 0 {
				downloadables = append(downloadables, *file)
			}
		}
		if _, valid := MyPeer.Comm["tracker"]; valid { // Supposing this tests if the tracker is connected.
			availableCommands = append(availableCommands, 1)
			fmt.Println("(1) Look for files")
			if len(tools.RemoteFiles) > 0 {
				availableCommands = append(availableCommands, 2)
				fmt.Println("(2) Get peers' infos")
			}
		}
		if len(downloadables) > 0 {
			availableCommands = append(availableCommands, 3)
			fmt.Println("(3) Download a file")
		}
		fmt.Println("(4) Handle a tracker or peer [Advanced]")
		fmt.Println("(5) Connect to a peer [Advanced]")
		fmt.Println("(6) Close connection to a tracker or peer [Advanced]")
		fmt.Println("(0) Exit")

		fmt.Print("\u001B[92mEnter a command:\u001B[39m ")
		command, err := strconv.Atoi(readInput())
		if err != nil {
			command = -1
		}
		switch command {
		case 1: // Look
			if !slices.Contains(availableCommands, 1) {
				continue
			}
			fmt.Println("\u001B[92mPerforming 'look' command...\u001B[39m")
			fmt.Print("Which criteria you are looking for: ") // What if wrong format ?
			commandInput := readInput()
			if commandInput == "" {
				peer.WriteReadConnection(MyPeer.Comm["tracker"], &MyPeer, "look []\n")
			} else {
				peer.WriteReadConnection(MyPeer.Comm["tracker"], &MyPeer, "look ["+commandInput+"]\n")
			}
		case 2: // Getfile
			if !slices.Contains(availableCommands, 2) {
				continue
			}
			fmt.Println("\u001B[92mPerforming 'getfile' command...\u001B[39m")
			var remoteFileKeys []string
			i := 0
			for key, file := range tools.RemoteFiles {
				remoteFileKeys = append(remoteFileKeys, key)
				fmt.Printf("(%d) %s %s\n", i, file.Name, key)
				i++
			}

			fmt.Print("\u001B[92mSelection:\u001B[39m ")
			input := readInput()
			fileNum, err := strconv.Atoi(input)
			if err != nil || fileNum < 0 || fileNum >= len(remoteFileKeys) {
				fmt.Println("\u001B[91mInvalid selection.\u001B[39m")
				continue
			}

			selectedFile := remoteFileKeys[fileNum]

			fmt.Printf("\u001B[92mPerforming 'getfile' command for file %s...\u001B[39m\n", selectedFile)
			peer.WriteReadConnection(MyPeer.Comm["tracker"], &MyPeer, fmt.Sprintf("getfile %s\n", selectedFile))

		case 3: // Download
			if !MyPeer.IsEmpty() { // Not sure if the test is necessary.
				if !slices.Contains(availableCommands, 3) {
					continue
				}
				//peer.WriteReadConnection(MyPeer.Comm["tracker"], &MyPeer, "look []\n") // ?????????
				fmt.Print("\u001B[92mHere all the files you can download :\u001B[39m \n")
				var remoteFileKeys []string
				i := 0
				for _, file := range downloadables {
					fmt.Printf("(%d) %s %s (%d)\n", i, file.Name, file.Key, file.Size)
					remoteFileKeys = append(remoteFileKeys, file.Key)
					i++
				}

				fmt.Print("\u001B[92mSelection:\u001B[39m ")
				input := readInput()
				fileNum, err := strconv.Atoi(input)
				if err != nil || fileNum < 0 || fileNum >= len(tools.RemoteFiles) {
					fmt.Println("\u001B[91mInvalid selection.\u001B[39m")
					continue
				}

				selectedFile := remoteFileKeys[fileNum]
				go MyPeer.Downloading(selectedFile)
				success := peer.WaitFor(peer.DlDone, true)
				if !success {
					fmt.Println("\u001B[91mDownload time out\u001B[39m")
				}
			}

		case 4: // Handle
			if !MyPeer.IsEmpty() {
				handlePeer(&MyPeer, "handle")
			} else { // Still needed ?
				fmt.Println("\u001B[91mYou need to launch a peer first.\u001B[39m")
			}
		case 5: // Connect
			if !MyPeer.IsEmpty() {
				fmt.Print("\u001B[92mEnter the IP and port of peer you want to connect to:\u001B[39m ")
				input := readInput()
				if len(input) == 2 {
					MyPeer.ConnectTo(string(input[0]), string(input[1]))
				} else if len(input) == 1 {
					MyPeer.ConnectTo("localhost", string(input[0]))
				} else {
					fmt.Println("\u001B[91mInvalid input.\u001B[39m")
				}
			} else { // Still needed ?
				fmt.Println("\u001B[92mYou need to launch a peer first.\u001B[39m")
			}
		case 6: // Close
			if !MyPeer.IsEmpty() {
				handlePeer(&MyPeer, "close")
			} else { // Still needed ?
				fmt.Println("\u001B[91mYou need to launch a peer first.\u001B[39m")
			}
		case 0: // Exit
			fmt.Println("Ending the program, closing all connections.")
			for key := range MyPeer.Comm {
				MyPeer.Close(key)
			}
			os.Exit(0)
		default:
			fmt.Println("\u001B[91mCommand not found. Available commands: lp (launch a peer), download, connect, handle, close, exit\u001B[39m")
		}

	}
}

func main() {
	peer.DlDone = make(chan struct{})
	peer.ResponsesRemainingUpdated = make(chan struct{})
	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			s := <-sigchnl
			for key := range MyPeer.Comm {
				MyPeer.Close(key)
			}
			fmt.Println("Received :", s, ". Please exit the code properly by typing exit :)")
			os.Exit(1)
		}
	}()
	tools.LogFile, _ = tools.OpenLog()
	tools.WriteLog("Lancement du peer...")
	inputProg()
	select {}
}
