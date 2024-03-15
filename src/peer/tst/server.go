package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func errorCheck(err error, str string) {
	if err != nil {
		fmt.Println("Error:", err, str)
		panic(err)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buffer := make([]byte, 256)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Reading error:", err)
			return
		}

		if n > 0 {
			fmt.Print("<Message> ", string(buffer[:n]))
			_, err := conn.Write(buffer[:n])
			if err != nil {
				fmt.Println("Writing error:", err)
				return
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "ERROR, no port provided\n")
		os.Exit(1)
	}

	portno, err := strconv.Atoi(os.Args[1])
	errorCheck(err, "arg error")

	l, err := net.Listen("tcp", ":"+strconv.Itoa(portno))
	errorCheck(err, "listening error")
	defer l.Close()

	fmt.Println("Server listening on port", portno)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Acceptation error:", err)
			continue
		}

		// Démarrez une goroutine pour gérer la connexion
		go handleConnection(conn)
	}
}

// func (p *Peer) StartListening() {
//
// 	// peers := make(map[string]Peer)
// 	server, err := net.Listen("tcp", p.IP+":"+p.Port)
// 	errorCheck(err)
//
// 	defer server.Close()
// 	fmt.Println("Hello I'm", p.Type)
// 	fmt.Println("Listening on " + p.IP + ":" + p.Port)
// 	fmt.Println("Waiting for client...\n")
//
// 	conn, err := server.Accept()
// 	errorCheck(err)
// 	defer conn.Close()
// 	if p.Type == "tracker" {
// 		//tools de kevin
// 		for {
// 			buffer := make([]byte, 256)
// 			n, err := conn.Read(buffer)
// 			errorCheck(err)
//
// 			switch string(buffer[:n]) {
// 			case "sendfile\n":
// 				fmt.Println("receiving the file :)", string(buffer))
// 				go p.receivefile("receivefile.txt", conn)
// 			}
// 			if n > 0 {
// 				fmt.Println("tracker :> ", string(buffer[:n]))
// 				_, err := conn.Write(buffer[:n])
// 				errorCheck(err)
// 			}
// 		}
// 	} else {
// 		fmt.Printf("Do you want to send or receive a file? (send/receive): from %s btw\n", p.Type)
// 		reader := bufio.NewReader(os.Stdin)
// 		option, _ := reader.ReadString('\n')
// 		option = strings.TrimSpace(option)
//
// 		switch option {
// 		case "send":
// 			go p.sendfile(conn)
// 		case "receive":
// 			go p.receivefile("./received.txt", conn)
// 		default:
// 			fmt.Println("Invalid option. Exiting.")
// 			os.Exit(1)
// 		}
// 	}
// }
//
// func (p *Peer) sendfile(conn net.Conn) {
// 	var file *os.File
// 	var err error
// 	if p.files == nil {
// 		file, err = os.Open("exemple.txt")
// 	} else {
// 		file, err = os.Open("exemple.txt")
// 	}
// 	errorCheck(err)
// 	defer file.Close()
//
// 	// Create a buffer to read the file in chunks
// 	buffer := make([]byte, 1024)
//
// 	for {
// 		// Read from the file into the buffer
// 		bytesRead, err := file.Read(buffer)
// 		if err == io.EOF {
// 			break
// 		}
// 		errorCheck(err)
//
// 		_, err = conn.Write(buffer[:bytesRead])
// 		errorCheck(err)
// 	}
// }
//
// func (p *Peer) receivefile(filename string, conn net.Conn) {
// 	defer conn.Close()
//
// 	file, err := os.Create(filename)
// 	errorCheck(err)
// 	defer file.Close()
//
// 	// Create a buffer to read the file data in chunks
// 	buffer := make([]byte, 1024)
//
// 	for {
// 		// Read from the connection into the buffer
// 		bytesRead, err := conn.Read(buffer)
// 		errorCheck(err)
// 		// fmt.Println(buffer, string(buffer[:bytesRead]))
// 		// Write the buffer data to the file
// 		_, err = file.Write(buffer[:bytesRead])
// 		errorCheck(err)
// 	}
// }
//
// func (p *Peer) StartPeer() {
// 	go p.StartListening()
// 	time.Sleep(2 * time.Second)
// 	tracker := GetConfig()
// 	serv_tcp_addr, err := net.ResolveTCPAddr("tcp", tracker.IP+":"+tracker.Port)
// 	errorCheck(err)
//
// 	sockfd, err := net.DialTCP("tcp", nil, serv_tcp_addr)
// 	errorCheck(err)
//
// 	reader := bufio.NewReader(os.Stdin)
// 	for {
// 		message, _ := reader.ReadString('\n')
// 		switch message {
// 		case "exit\n":
// 			os.Exit(1)
// 		case "sendfile\n":
// 			fmt.Println("Sending waiting for reception...")
// 			go p.sendfile(sockfd)
// 		}
// 		sockfd.Write([]byte(message))
//
// 		buffer := make([]byte, 256)
// 		n, err := sockfd.Read(buffer)
// 		errorCheck(err)
// 		fmt.Println(string(buffer[:n]))
// 	}
// }
//
// // Reprise du téléchargement en cas d'erreur avec un fichier d'opérations.
