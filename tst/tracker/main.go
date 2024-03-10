package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

func errorFunc(err error) {
	fmt.Println("Error:", err)
	os.Exit(1)
}

func main() {
	var sockfd net.Conn
	var portno int
	var n int
	var err error

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage %s hostname port\n", os.Args[0])
		os.Exit(1)
	}

	portno, err = strconv.Atoi(os.Args[2])
	if err != nil {
		errorFunc(err)
	}

	server := os.Args[1]
	serv_tcp_addr, err := net.ResolveTCPAddr("tcp", server+":"+strconv.Itoa(portno))
	if err != nil {
		errorFunc(err)
	}

	sockfd, err = net.DialTCP("tcp", nil, serv_tcp_addr)
	if err != nil {
		errorFunc(err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		message, _ := reader.ReadString('\n')
		sockfd.Write([]byte(message))

		buffer := make([]byte, 256)
		n, err = sockfd.Read(buffer)
		if err != nil {
			errorFunc(err)
		}
		fmt.Print(string(buffer[:n]))
	}
}