package peer_package

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"path/filepath"
	"peerproject/tools"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sys/unix"
)

func (p *Peer) Close(t string) {
	fmt.Printf("[\u001B[0;33m%s\u001B[39m] Connection closed\n", t)
	tools.WriteLog("[%s] Connection closed\n", t)
	p.Comm[t].Close()
	delete(p.Comm, t)
}

type Job struct {
	conn net.Conn
	data string
}

var attempts = struct {
	sync.RWMutex
	m map[net.Conn]int
}{m: make(map[net.Conn]int)}

func worker(jobs chan Job, p *Peer) {
	for job := range jobs {
		data := job.data
		conn := job.conn

		attempts.RLock()
		attempt, ok := attempts.m[conn]
		attempts.RUnlock()

		if !ok {
			fmt.Println("Connection not found")
			continue
		}

		mess := data
		mess = strings.Split(mess, "\n")[0]
		input := strings.Split(mess, " ")[0]

		fmt.Printf("[\u001B[0;33m%s\u001B[39m]: %s\n", conn.RemoteAddr().String(), mess)
		tools.WriteLog("[%s]: %s\n", conn.RemoteAddr().String(), mess)

		switch input {
		case "interested":
			valid, data := tools.InterestedCheck(mess)
			if valid {
				file := p.Files[data.Key]
				buff := "have " + data.Key + " " + tools.BufferMapToString(*file.Peers["self"].BufferMaps[data.Key]) + "\n"
				fmt.Printf("\033[0;35mSending to\u001B[39m [\u001B[0;33m%s\u001B[39m]:%s\n", conn.RemoteAddr().String(), buff)
				tools.WriteLog("Sending to %s:%s\n", conn.RemoteAddr().String(), buff)
				_, err := conn.Write([]byte(buff))
				errorCheck(err)
				break
			} else { // TODO : En faire une fonction
				attempts.Lock()
				attempts.m[conn]--
				attempt--
				attempts.Unlock()
				checkAttempts(conn, attempt)
			}
		case "getpieces":
			fmt.Println(conn.RemoteAddr().String(), ":", mess)
			valid, data := tools.GetPiecesCheck(mess)
			if valid {
				fdf, err := os.OpenFile(filepath.Join(tools.GetValueFromConfig("Peer", "path"), "/"+p.Files[data.Key].Name), os.O_CREATE|os.O_RDWR, os.FileMode(0777))
				errorCheck(err)
				response := "data " + data.Key + " ["
				for _, piece := range data.Pieces {
					_, err := fdf.Seek(int64(piece*p.Files[data.Key].PieceSize), 0)
					errorCheck(err)
					// tempBuff := make([]byte, p.Files[data.Key].PieceSize)
					tempBuff := tools.BufferMap{Length: p.Files[data.Key].PieceSize * 8, BitSequence: make([]byte, p.Files[data.Key].PieceSize*8)}
					fdf.Read(tempBuff.BitSequence)
					response += strconv.Itoa(piece) + ":" + tools.BufferMapToString(tempBuff) + " "
					// tempBuff := make([]byte, p.Files[data.Key].PieceSize)
					// fdf.Read(tempBuff)
					// response += strconv.Itoa(piece) + ":" + string(tempBuff) + " "
				}
				response = strings.TrimSuffix(response, " ")
				response += "]\n"
				fmt.Println(conn.LocalAddr().String(), ":", response)
				_, err = conn.Write([]byte(response))
				errorCheck(err)
				break
			} else {
				attempts.Lock()
				attempts.m[conn]--
				attempt--
				attempts.Unlock()
				checkAttempts(conn, attempt)
			}
			// TODO : Faire en sorte que ça s'envoie tous les 3 dl de pieces.
		case "have":
			valid, data := tools.HaveCheck(mess)
			if valid {
				response := "have " + data.Key + " " + tools.BufferMapToString(*p.Files[data.Key].Peers[conn.LocalAddr().String()].BufferMaps[data.Key]) + "\n"
				_, err := conn.Write([]byte(response))
				errorCheck(err)
			} else {
				attempts.Lock()
				attempts.m[conn]--
				attempt--
				attempts.Unlock()
				checkAttempts(conn, attempt)
			}
		case "exit":
			conn.Close()
			return
		default:
			attempts.Lock()
			attempts.m[conn]--
			attempt--
			attempts.Unlock()
			checkAttempts(conn, attempt)
		}
	}
}

func checkAttempts(conn net.Conn, attempt int) {
	buffer := "Invalid command you have " + strconv.Itoa(attempts.m[conn]) + " tries remaining\n"
	fmt.Println(conn.LocalAddr(), buffer)
	fmt.Printf("[\u001B[0;33m%s\u001B[39m] Invalid command. %d tries remaining.\n", conn.RemoteAddr().String(), attempts.m[conn])
	tools.WriteLog("[%s] Invalid command. %d tries remaining.\n", conn.RemoteAddr().String(), attempts.m[conn])
	_, err := conn.Write([]byte(buffer))
	errorCheck(err)

	if attempt <= 0 {
		buffer := "Invalid command you have no tries remaining, connection is closed...\n"
		fmt.Printf("[\u001B[0;33m%s\u001B[39m] Connection closed due to invalid messages.\n", conn.RemoteAddr().String())
		tools.WriteLog("[%s] Connection closed due to invalid messages.\n", conn.RemoteAddr().String())
		_, err := conn.Write([]byte(buffer)) // ?? Pas besoin de prévenir, la connection est fermée à la ligne d'après.
		conn.Close()
		errorCheck(err)
	}
}

func (p *Peer) startListening() { // You are stuck here if the IP is not valid.
	l, err := net.Listen("tcp", p.IP+":"+p.Port)
	for err != nil {
		fmt.Println("Listen error:", err)
		fmt.Println("Please choose an other port this one is already bind :")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		l, err = net.Listen("tcp", p.IP+":"+text)
		p.Port = text
	}
	defer l.Close()

	fmt.Println("Server listening on port", p.Port)
	tools.WriteLog("Server listening on port: %s\n", p.Port)

	epfd, err := unix.EpollCreate1(0)
	if err != nil {
		fmt.Println("EpollCreate1 error:", err)
		return
	}
	defer unix.Close(epfd)

	var events [128]unix.EpollEvent

	max_concurrence, err := strconv.Atoi(tools.GetValueFromConfig("Peer", "max_concurrency"))
	errorCheck(err)
	max_attempts, err := strconv.Atoi(tools.GetValueFromConfig("Peer", "max_message_attempts"))
	errorCheck(err)
	max_peers, err := strconv.Atoi(tools.GetValueFromConfig("Peer", "max_peers"))
	errorCheck(err)

	jobs := make(chan Job, max_peers)
	for i := 0; i < max_concurrence; i++ {
		go worker(jobs, p)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		tcpConn, ok := conn.(*net.TCPConn)
		if !ok {
			log.Fatal("Not a TCP connection")
		}

		file, err := tcpConn.File()
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}

		fd := int(file.Fd())
		event := unix.EpollEvent{
			Events: unix.EPOLLIN,
			Fd:     int32(fd),
		}
		if err := unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, fd, &event); err != nil {
			fmt.Println("EpollCtl error:", err)
			conn.Close()
			continue
		}
		attempts.Lock()
		attempts.m[conn] = max_attempts
		attempts.Unlock()

		go func() {
			defer conn.Close()
			defer unix.EpollCtl(epfd, unix.EPOLL_CTL_DEL, fd, nil)

			for {
				nevents, err := unix.EpollWait(epfd, events[:], -1)
				if err != nil {
					fmt.Println("EpollWait error:", err)
					return
				}

				for ev := 0; ev < nevents; ev++ {
					if (events[ev].Events&unix.EPOLLHUP) != 0 || (events[ev].Events&unix.EPOLLERR) != 0 {
						return
					} else if (events[ev].Events & unix.EPOLLIN) != 0 {
						var data string
						buf := make([]byte, 32768)
						var err error = nil
						var n int
						for len(data) == 0 || data[len(data)-1] != '\n' || err != nil {
							conn.SetReadDeadline(time.Now().Add(time.Duration(float64(7) * math.Pow(10, 9))))
							fd, err = conn.Read(buf)
							n += fd
							// errorCheck(nerr)
							data += string(buf[:fd])
						}
						conn.SetReadDeadline(time.Time{})
						fmt.Println(string(data))
						jobs <- Job{conn, data}
					}
					// 	data := make([]byte, 32768)
					// 	_, err := conn.Read(data)
					// 	if err != nil {
					// 		fmt.Println("Read error:", err)
					// 		return
					// 	}
					//
					// 	jobs <- Job{conn, data}
					// }
				}
			}
		}()
	}
}
