package peer_package

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"peerproject/tools"
	"strconv"
	"strings"
	"time"
	// "peerproject/tools"
)

func (p *Peer) HelloTrack(t Peer) {
	message := "announce listen " + p.Port + " seed ["
	for _, valeur := range p.Files {
		name, size, pieceSize, key, isEmpty := valeur.GetFile()
		if isEmpty {
			message += fmt.Sprintf(`%s %d %d %s `, name, size, pieceSize, key)
		} else {
			break
		}
	}
	message = strings.TrimSuffix(message, " ")
	message += "]\n"
	conn, err := net.Dial("tcp", t.IP+":"+t.Port)
	errorCheck(err)
	// defer conn.Close()
	message = string(message)
	_, err = conn.Write([]byte(message))
	errorCheck(err)
	buffer := make([]byte, 32768)
	n, err := conn.Read(buffer)
	errorCheck(err)
	fmt.Print("< ", string(buffer[:n]))
	p.Comm[conn.RemoteAddr().String()] = conn
}

func (p *Peer) sendupdate(t Peer) {
	var seed string
	var leech string
	for {
		seed = ""
		leech = ""
		updateValue, _ := strconv.Atoi(tools.GetValueFromConfig("Peer", "update_time"))
		time.Sleep(time.Duration(float64(updateValue) * math.Pow(10, 9)))
		message := "update seed ["
		for _, valeur := range p.Files {
			key, BitSequence, notEmpty := valeur.GetFileUpdate()
			if notEmpty {
				temp := 0
				// for k := range len(BitSequence) {
				for k := 0; k < len(BitSequence); k++ {
					if !tools.ByteArrayCheck(BitSequence, k) {
						temp++
					}
				}
				if temp > 0 {
					leech += fmt.Sprintf(`%s `, key)
				} else {
					seed += fmt.Sprintf(`%s `, key)
				}
			} else {
				break
			}
		}
		seed = strings.TrimSuffix(seed, " ")
		leech = strings.TrimSuffix(leech, " ")
		message += seed + "] leech [" + leech + "]\n"
		conn, err := net.Dial("tcp", t.IP+":"+t.Port)
		errorCheck(err)
		_, err = conn.Write([]byte(message))
		errorCheck(err)
	}
}

// TODO : Remote file stockage lors d une demande au tracker
// TODO : Faire les chanegement dans les fonctions car changement de []file en map.
// TODO : Finir la gestion des messages notamment le dl et du buffermap avec les fonctions de tools.
func WriteReadConnection(conn net.Conn, p *Peer, mess ...string) {
	var message string
	if len(mess) == 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\u001B[92mWaiting for input... :\u001B[39m")
		message, _ = reader.ReadString('\n')

	} else {
		message = mess[0]
	}
	fmt.Println(conn.LocalAddr().String(), ":", message)
	if message == "exit\n" {
		os.Exit(1)
	}
	conn.Write([]byte(message))

	buffer := make([]byte, 32768)
	fd, nerr := conn.Read(buffer)
	errorCheck(nerr)
	if fd > 0 {
		mess := string(buffer[:fd])
		mess = strings.TrimSuffix(mess, "\n")
		input := strings.Split(mess, " ")[0]
		switch input {
		case "data", "data\n":
			fmt.Printf("%s\n", mess)
			valid, data := tools.DataCheck(mess)
			if valid {
				fmt.Println(conn.RemoteAddr(), ":", mess)
				// fmt.Println("ET OUAIS .....")
				os.MkdirAll(filepath.Join("./", tools.GetValueFromConfig("Peer", "path"), "/", "tabernak"), os.FileMode(0777))
				file := p.Files[data.Key]
				// if !ok {
				// 	rFile := tools.RemoteFiles[data.Key]
				// 	p.Files[data.Key] = &tools.File{Name: rFile.Name, Size: rFile.Size, PieceSize: rFile.PieceSize, Key: data.Key}
				// 	tools.InitBufferMap(p.Files[data.Key])
				// 	file = p.Files[data.Key]
				// }
				fdf, err := os.OpenFile(filepath.Join(tools.GetValueFromConfig("Peer", "path"), "tabernak/file"), os.O_CREATE|os.O_RDWR, os.FileMode(0777))
				errorCheck(err)
				fdc, err := os.OpenFile(filepath.Join(tools.GetValueFromConfig("Peer", "path"), "tabernak/log"), os.O_CREATE|os.O_RDWR|os.O_APPEND, os.FileMode(0777))
				errorCheck(err)
				for i := 0; len(data.Pieces) > i; i++ {
					_, err := fdf.Seek(int64(data.Pieces[i].Index*file.PieceSize), 0)
					errorCheck(err)
					n, err := fdf.Write(data.Pieces[i].Data.BitSequence)
					if n <= 0 {
						fmt.Println("File as not being written. :(", err, data.Pieces[i].Data.BitSequence)

					}
					_, err = fdc.WriteString("\n" + time.Now().String() + "Downloading the " + fmt.Sprint(data.Pieces[i].Index) + " piece.")
					errorCheck(err)
					tools.ByteArrayWrite(&file.Peers[conn.LocalAddr().String()].BufferMaps[data.Key].BitSequence, data.Pieces[i].Index)

				}
				fdf.Close()
				fdc.Close()
			} else {
				fmt.Println("\u001B[92mInvalid data response...\u001B[39m")
			}
		case "have", "have\n":
			valid, data := tools.HaveCheck(mess)

			if valid {
				tools.BufferMapCopy(&data.BufferMap, tools.RemoteFiles[data.Key].Peers[conn.RemoteAddr().String()].BufferMaps[data.Key]) // TODO Debug

				go p.progression(data.Key, p.Files[data.Key].Peers[conn.LocalAddr().String()].BufferMaps[data.Key].Length, conn)
				fmt.Println(conn.RemoteAddr(), ":", mess)
			} else {
				fmt.Println("\u001B[92mInvalid have response...\u001B[39m")
			}
		case "OK", "OK\n":
			fmt.Println(conn.RemoteAddr(), ":", mess)
		case "list", "list\n":
			valid, _ := tools.ListCheck(mess)
			if valid {
				fmt.Println(conn.RemoteAddr(), ":", mess)
			} else {
				fmt.Println("\u001B[92mInvalid list response...\u001B[39m")
			}
		case "peers", "peers\n":
			valid := tools.PeersCheck(mess)
			if valid {
				fmt.Println(conn.RemoteAddr(), ":", mess)
				p.interested(string(input[1]))
			} else {
				fmt.Println("\u001B[92mInvalid peers response...\u001B[39m")
			}
		case "\u001B[92mInvalid command you have no tries remaining, connection is closed...\u001B[39m":
			p.Close(conn.RemoteAddr().String())
		default:
			// panic("valeur par default et pas parmi la liste")
			// je dois prendre en compte que si je n'ai plus d'essai de fermer de mon cote la conn.
			fmt.Println(string(buffer[:]))

		}
	}
}

func (p *Peer) interested(key string) {
	temp := tools.RemoteFiles[key]
	l := len(temp.Peers)
	for max, _ := strconv.Atoi(tools.GetValueFromConfig("Peer", "max_peers_to_connect")); max != 0; max-- {
		//random := rand.Intn(l)
		var randomPeer tools.Peer
		k := rand.Intn(l)
		for _, peer := range temp.Peers {
			if k == 0 {
				randomPeer = *peer
			}
			k--
		}
		p.ConnectTo(randomPeer.IP, fmt.Sprint(randomPeer.Port), "interested"+" "+key)
	}
}

func (p *Peer) progression(key string, length int, conn net.Conn) {
	i, _ := strconv.Atoi(tools.GetValueFromConfig("Peer", "progress_value"))
	for {
		if file, valid := p.Files["self"]; valid {
			if file.Peers[conn.LocalAddr().String()].BufferMaps[key].Length == length+i {
				WriteReadConnection(conn, p, "have "+key+" "+tools.BufferMapToString(*file.Peers[conn.LocalAddr().String()].BufferMaps[key]))
				i = 0
			}
		}
	}
}
func (p *Peer) ConnectTo(IP string, Port string, mess ...string) {
	conn, err := net.Dial("tcp", IP+":"+Port)
	errorCheck(err)
	// defer conn.Close()
	p.Comm[conn.RemoteAddr().String()] = conn
	fmt.Println(conn.LocalAddr(), " is connected to ", conn.RemoteAddr())
	if len(mess) == 0 {
		WriteReadConnection(conn, p)
	} else {
		WriteReadConnection(conn, p, mess...)
	}

	// go handleConnection(conn)
	//Handle the response here !
	// fmt.Print(string(buffer[:n]))
}
