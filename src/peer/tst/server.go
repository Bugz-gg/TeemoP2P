package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func errorCheck(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "ERROR, no port provided\n")
		os.Exit(1)
	}

	portno, err := strconv.Atoi(os.Args[1])
	errorCheck(err)

	l, err := net.Listen("tcp", ":"+strconv.Itoa(portno))
	errorCheck(err)
	defer l.Close()

	fmt.Println("Server listening on port", portno)

	conn, err := l.Accept()
	errorCheck(err)
	defer conn.Close()

	for {
		buffer := make([]byte, 256)
		n, err := conn.Read(buffer)
		errorCheck(err)

		if n > 0 {
			fmt.Print("<Message> ", string(buffer[:n]))
			_, err := conn.Write(buffer[:n])
			errorCheck(err)
		}
	}
}
