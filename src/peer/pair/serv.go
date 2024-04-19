package peer_package

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"peerproject/tools"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/sys/unix"
)

func (p *Peer) Close(t string) {
	fmt.Println("Carefull the connection with", t, "is now closed")
	p.Comm[t].Close()
	delete(p.Comm, t)
}

type Job struct {
	conn net.Conn
	data []byte
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

		mess := string(data[:])
		mess = strings.Split(mess, "\n")[0]
		input := strings.Split(mess, " ")[0]

		switch input {
		case "interested":
			valid, data := tools.InterestedCheck(mess)
			fmt.Println(valid)
			if valid {
				// fmt.Println("in if")
				file := p.Files[data.Key]
				buff := "have " + data.Key + " " + tools.BufferMapToString(*file.Peers["self"].BufferMaps[data.Key]) + "\n"
				fmt.Print(conn.LocalAddr(), buff)
				_, err := conn.Write([]byte(buff))
				errorCheck(err)
				break
			} else { // TODO : En faire une fonction
				attempts.Lock()
				attempts.m[conn]--
				attempt--
				attempts.Unlock()
				if attempt == 0 {
					buffer := "Invalid command you have no tries remaining, connection is closed..."
					fmt.Println(conn.LocalAddr(), buffer)
					_, err := conn.Write([]byte(buffer))
					conn.Close()
					errorCheck(err)

				} else {
					buffer := "Invalid command you have " + strconv.Itoa(attempts.m[conn]) + " tries remaining"
					fmt.Println(conn.LocalAddr(), buffer)
					_, err := conn.Write([]byte(buffer))
					errorCheck(err)

				}
			}
		case "getpieces", "getpieces\n":
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
				response += "]"
				fmt.Println(conn.LocalAddr().String(), ":", response)
				_, err = conn.Write([]byte(response))
				errorCheck(err)
				break
			} else {
				attempts.Lock()
				attempts.m[conn]--
				attempt--
				attempts.Unlock()
				if attempt == 0 {
					buffer := "Invalid command you have no tries remaining, connection is closed..."
					fmt.Println(conn.LocalAddr(), buffer)
					_, err := conn.Write([]byte(buffer))
					conn.Close()
					errorCheck(err)

				} else {
					buffer := "Invalid command you have " + strconv.Itoa(attempts.m[conn]) + " tries remaining"
					fmt.Println(conn.LocalAddr(), buffer)
					_, err := conn.Write([]byte(buffer))
					errorCheck(err)

				}
			}
			// TODO : Faire en sorte que ça s" envoie toutes les 3 dl de pieces.
		case "have", "have\n":
			valid, data := tools.HaveCheck(mess)
			if valid {
				response := "have " + data.Key + " " + tools.BufferMapToString(*p.Files[data.Key].Peers[conn.LocalAddr().String()].BufferMaps[data.Key])
				_, err := conn.Write([]byte(response))
				errorCheck(err)
			} else {
				attempts.Lock()
				attempts.m[conn]--
				attempt--
				attempts.Unlock()
				if attempt == 0 {
					buffer := "Invalid command you have no tries remaining, connection is closed..."
					fmt.Println(conn.LocalAddr(), buffer)
					_, err := conn.Write([]byte(buffer))
					conn.Close()
					errorCheck(err)

				} else {
					buffer := "Invalid command you have " + strconv.Itoa(attempts.m[conn]) + " tries remaining"
					fmt.Println(conn.LocalAddr(), buffer)
					_, err := conn.Write([]byte(buffer))
					errorCheck(err)
				}
			}
		case "exit", "exit\n":
			conn.Close()
			return
		default:
			attempts.Lock()
			attempts.m[conn]--
			attempt--
			attempts.Unlock()
			if attempt == 0 {
				buffer := "Invalid command you have no tries remaining, connection is closed..."
				fmt.Println(conn.LocalAddr(), buffer)
				_, err := conn.Write([]byte(buffer))
				conn.Close()
				errorCheck(err)

			} else {
				buffer := "Invalid command you have " + strconv.Itoa(attempts.m[conn]) + " tries remaining"
				fmt.Println(conn.LocalAddr(), buffer)
				_, err := conn.Write([]byte(buffer))
				errorCheck(err)

			}
		}
	}
}

func (p *Peer) startListening() {
	l, err := net.Listen("tcp", p.IP+":"+p.Port)
	if err != nil {
		fmt.Println("Listen error:", err)
		return
	}
	defer l.Close()

	fmt.Println("Server listening on port", p.Port)

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
						data := make([]byte, 0, 32768)
						buf := make([]byte, 32768)
						for {
							n, err := conn.Read(buf)
							errorCheck(err)
							data = append(data, buf[:n]...)

							if bytes.Contains(data, []byte{'\n'}) {
								break
							}
						}
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
