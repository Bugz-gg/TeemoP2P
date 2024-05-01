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
	"sync"
	"time"
)

var mutex sync.Mutex
var previousMessage string // To know if it was an interested or a have.<ScrollWheelDown>
var rare bool

func (p *Peer) HelloTrack(t Peer) {
	timeout, _ := strconv.Atoi(tools.GetValueFromConfig("Peer", "timeout"))
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
	buffer := make([]byte, 1024)
	err = conn.SetReadDeadline(time.Now().Add(time.Duration(float64(timeout) * math.Pow(10, 9))))
	errorCheck(err)
	n, err := conn.Read(buffer)
	errorCheck(err)
	err = conn.SetReadDeadline(time.Time{})
	if err == nil {
		fmt.Print(string(buffer[:n]))
	}
	p.Comm[conn.RemoteAddr().String()] = conn
	p.Comm["tracker"] = p.Comm[conn.RemoteAddr().String()]
}

func (p *Peer) sendupdate(t Peer) {
	var seed string
	var leech string
	timeout, _ := strconv.Atoi(tools.GetValueFromConfig("Peer", "timeout"))
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
		buffer := make([]byte, 1024)
		err = conn.SetReadDeadline(time.Now().Add(time.Duration(float64(timeout) * math.Pow(10, 9))))
		n, err := conn.Read(buffer)
		err = conn.SetReadDeadline(time.Time{})
		if err == nil {
			fmt.Print(string(buffer[:n]))
		}
	}
}

func (p *Peer) rarepiece() {
	time.Sleep(time.Second)
	conn := p.Comm["tracker"]
	t, _ := strconv.Atoi(tools.GetValueFromConfig("Peer", "time_dl_rare_piece"))
	var byteArray []int
	var connArray map[int][]net.Conn
	for {
		WriteReadConnection(conn, p, "look []\n")
		for key := range tools.RemoteFiles {
			if _, valid := p.Files[key]; !valid || !tools.ArrayCheck(p.Files[key].Peers["self"].BufferMaps[key].BitSequence) {
				k := max(int(tools.BufferBitSize(*tools.RemoteFiles[key])/4), 1)
				// fmt.Println(k, tools.RemoteFiles[key].PieceSize, tools.RemoteFiles[key].Size)
				rareArray := make([]int, k)
				mutex.Lock()
				rare = true
				mutex.Unlock()
				WriteReadConnection(conn, p, "getfile "+key+"\n")
				for connRem := range tools.RemoteFiles[key].Peers {
					p.ConnectTo(tools.RemoteFiles[key].Peers[connRem].IP, tools.RemoteFiles[key].Peers[connRem].Port, "interested "+key+"\n")
					j := 0
					if byteArray == nil {
						byteArray = make([]int, tools.RemoteFiles[key].Peers[connRem].BufferMaps[key].Length)
						connArray = make(map[int][]net.Conn, tools.RemoteFiles[key].Peers[connRem].BufferMaps[key].Length)
					}
					for i := range tools.RemoteFiles[key].Peers[connRem].BufferMaps[key].Length {
						// fmt.Print(i, tools.RemoteFiles[key].Peers[connRem].BufferMaps[key].BitSequence, byteArray, connArray, tools.RemoteFiles[key].Peers[connRem].BufferMaps[key].Length)
						if byteArray[j] != int(math.Inf(1)) {
							if tools.ByteArrayCheck(tools.RemoteFiles[key].Peers[connRem].BufferMaps[key].BitSequence, i) {
								connArray[j] = append(connArray[j], p.Comm[connRem])
								byteArray[j] += 1
							} else if tools.ByteArrayCheck(p.Files[key].Peers["self"].BufferMaps[key].BitSequence, i) {
								byteArray[j] = int(math.Inf(1))
							}
						}
						j++
					}

					// fmt.Print(tools.RemoteFiles[key].Peers[connRem].BufferMaps[key].BitSequence, byteArray, connArray, tools.RemoteFiles[key].Peers[connRem].BufferMaps[key].Length)
				}
				minIndex := 0
				for i := 0; i < k; i++ {
					for x := i + 1; x < len(byteArray); x++ {
						if byteArray[x] < byteArray[minIndex] {
							minIndex = x
						}
					}
					// fmt.Println(i, minIndex)
					byteArray[i], byteArray[minIndex] = byteArray[minIndex], byteArray[i]
					rareArray[i] = minIndex
				}
				for index := range rareArray { // TODO : Can improve is in case its only one peer to send it only once
					fmt.Println(index, rareArray)
					// fmt.Println(connArray, connArray[index], len(connArray[index]), rareArray, byteArray)
					WriteReadConnection(connArray[index][rand.Intn(len(connArray[index]))], p, "getpieces "+key+" ["+strconv.Itoa(index)+"]\n") // Beter with sprintf maybe
					time.Sleep(time.Millisecond)
				}
				break
			}
		}
		time.Sleep(time.Duration(math.Pow(10, 9) * float64(t)))
	}

}

// func (p *Peer) Downloading(key string) {
// 	time.Sleep(time.Second)
// 	conn := p.Comm["tracker"]
// 	var byteArray []int
// 	var dontHave []int
// 	var connArray map[int][]net.Conn
// 	for index := range p.Files[key].Peers["self"].BufferMaps[key].Length {
// 		if !tools.ByteArrayCheck(p.Files[key].Peers["self"].BufferMaps[key].BitSequence, index) {
// 			dontHave = append(dontHave, index)
// 		}
// 	}
// 	byteArray = make([]int, len(dontHave))
// 	connArray = make(map[int][]net.Conn, len(dontHave))
// 	WriteReadConnection(p.Comm["tracker"], p, "look [key="+key+"]\n")
//
// 	for connRem := range tools.RemoteFiles[key].Peers {
// 		p.ConnectTo(tools.RemoteFiles[key].Peers[connRem].IP, tools.RemoteFiles[key].Peers[connRem].Port, "interested "+key+"\n")
// 		j := 0
// 		for i := range dontHave {
// 			if tools.ByteArrayCheck(tools.RemoteFiles[key].Peers[connRem].BufferMaps[key].BitSequence, i) {
// 				connArray[i] = append(connArray[i], p.Comm[connRem])
// 				byteArray[j] += 1
// 			}
// 			j++
// 		}
//   }
//
// 		for tools.BitCount(p.Files[key].Peers["self"].BufferMaps[key]) == nbPieces {
// 			k := max(int(tools.BufferBitSize(*tools.RemoteFiles[key])/4), 1)
// 			// fmt.Println(k, tools.RemoteFiles[key].PieceSize, tools.RemoteFiles[key].Size)
// 			rareArray := make([]int, k)
// 			mutex.Lock()
// 			rare = true
// 			mutex.Unlock()
// 			WriteReadConnection(conn, p, "getfile "+key+"\n")
//
// 			// fmt.Print(tools.RemoteFiles[key].Peers[connRem].BufferMaps[key].BitSequence, byteArray, connArray, tools.RemoteFiles[key].Peers[connRem].BufferMaps[key].Length)
// 		}
// 		minIndex := 0
// 		for i := 0; i < k; i++ {
// 			for x := i + 1; x < len(byteArray); x++ {
// 				if byteArray[x] < byteArray[minIndex] {
// 					minIndex = x
// 				}
// 			}
// 			// fmt.Println(i, minIndex)
// 			byteArray[i], byteArray[minIndex] = byteArray[minIndex], byteArray[i]
// 			rareArray[i] = minIndex
// 		}
// 		for index := range rareArray { // TODO : Can improve is in case its only one peer to send it only once
// 			fmt.Println(index, rareArray)
// 			// fmt.Println(connArray, connArray[index], len(connArray[index]), rareArray, byteArray)
// 			WriteReadConnection(connArray[index][rand.Intn(len(connArray[index]))], p, "getpieces "+key+" ["+strconv.Itoa(index)+"]\n") // Beter with sprintf maybe
// 			time.Sleep(time.Millisecond)
// 		}
// 	}
//
// }

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
	var eom string
	var n int
	var nerr error = nil
	var fd int
	for len(eom) == 0 || eom[len(eom)-1] != '\n' || nerr != nil {
		conn.SetReadDeadline(time.Now().Add(time.Duration(float64(7) * math.Pow(10, 9))))
		fd, nerr = conn.Read(buffer)
		n += fd
		// errorCheck(nerr)
		eom += string(buffer[:fd])
	}
	conn.SetReadDeadline(time.Time{})
	if n > 0 {
		mess := eom
		mess = strings.TrimSuffix(mess, "\n")
		input := strings.Split(mess, " ")[0]
		switch input {
		case "data", "data\n":
			fmt.Printf("%s\n", mess)
			valid, data := tools.DataCheck(mess)
			if valid {
				fmt.Println(conn.RemoteAddr(), ":", mess)
				path := tools.GetValueFromConfig("Peer", "path")
				if path == "" {
					path = "share"
				}
				os.MkdirAll(filepath.Join("./", path, "/", p.Files[data.Key].Name), os.FileMode(0777))
				fmt.Println("the file is", p.Files[data.Key])
				file := p.Files[data.Key]
				fdf, err := os.OpenFile(filepath.Join(path, p.Files[data.Key].Name+"/file"), os.O_CREATE|os.O_RDWR, os.FileMode(0777))
				errorCheck(err)
				fdc, err := os.OpenFile(filepath.Join(path, p.Files[data.Key].Name+"/manifest"), os.O_CREATE|os.O_RDWR|os.O_APPEND, os.FileMode(0777))
				errorCheck(err)
				for i := 0; len(data.Pieces) > i; i++ {
					_, err := fdf.Seek(int64(data.Pieces[i].Index*file.PieceSize), 0)
					errorCheck(err)
					var n int
					if file.Size%file.PieceSize != 0 && data.Pieces[i].Index+1 == tools.BufferBitSize(*file) {
						n, err = fdf.Write(data.Pieces[i].Data.BitSequence[:file.Size%file.PieceSize])
						errorCheck(err)

					} else {
						n, err = fdf.Write(data.Pieces[i].Data.BitSequence)
						errorCheck(err)
					}
					if n <= 0 {
						fmt.Println("File as not being written. :(", err, data.Pieces[i].Data.BitSequence)
					}
					_, err = fdc.WriteString("\n" + time.Now().String() + "Downloading the " + fmt.Sprint(data.Pieces[i].Index) + " piece.")
					errorCheck(err)
					tools.ByteArrayWrite(&file.Peers["self"].BufferMaps[data.Key].BitSequence, data.Pieces[i].Index)
				}
				if hash := tools.GetMD5Hash(filepath.Join(path, p.Files[data.Key].Name+"/file")); hash == data.Key {
					err = os.Rename(filepath.Join(path, p.Files[data.Key].Name+"/file"), filepath.Join(path, p.Files[data.Key].Name))
					errorCheck(err)

					err = os.RemoveAll(filepath.Join(path, p.Files[data.Key].Name))
					errorCheck(err)
				}

				fdf.Close()
				fdc.Close()
			} else {
				fmt.Println("\u001B[92mInvalid data response...\u001B[39m")
			}
		case "have", "have\n":
			valid, data := tools.HaveCheck(mess)

			if valid {
				peer := tools.RemoteFiles[data.Key].Peers[conn.RemoteAddr().String()]
				if peer.BufferMaps == nil {
					peer.BufferMaps = make(map[string]*tools.BufferMap)
				}
				bufferMap := peer.BufferMaps[data.Key]
				tools.BufferMapCopy(&bufferMap, &data.BufferMap)
				if peer.BufferMaps == nil {
					peer.BufferMaps = make(map[string]*tools.BufferMap)
				}
				peer.BufferMaps[data.Key] = bufferMap

				if _, valid := p.Files[data.Key]; !valid {
					fil := tools.File{
						Name:      tools.RemoteFiles[data.Key].Name,
						Size:      tools.RemoteFiles[data.Key].Size,
						PieceSize: tools.RemoteFiles[data.Key].PieceSize,
						Key:       data.Key,
					}

					bufferMaps := make(map[string]*tools.BufferMap)
					buffermap := tools.InitBufferMap(tools.RemoteFiles[data.Key].Size, tools.RemoteFiles[data.Key].PieceSize)
					// fmt.Println(buffermap)
					bufferMaps[data.Key] = &buffermap
					//InitBufferMap(&fil)
					if fil.Peers == nil {
						fil.Peers = make(map[string]*tools.Peer)
					}
					fil.Peers["self"] = &tools.Peer{
						IP:         p.IP,
						Port:       p.Port,
						BufferMaps: bufferMaps,
					}
					if _, valid := fil.Peers["self"].BufferMaps[data.Key]; !valid {
						p.Files[data.Key].Peers["self"].BufferMaps = make(map[string]*tools.BufferMap)
					}
					fil.Peers[conn.LocalAddr().String()] = fil.Peers["self"]
					p.Files[data.Key] = &fil
				}
				if previousMessage == "interested" {
					// go p.progression(data.Key, conn)
					time.Sleep(2)
				} else {
					time.Sleep(2) // Handle new pieces.
				}

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
				key := strings.Split(mess, " ")[1]
				mutex.Lock()
				previousMessage = "interested"
				mutex.Unlock()
				if !rare {
					p.interested(key)
				}
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
	it := l
	for max, _ := strconv.Atoi(tools.GetValueFromConfig("Peer", "max_peers_to_connect")); max != 0 && it > 0; max-- {
		//random := rand.Intn(l)
		var randomPeer tools.Peer
		k := rand.Intn(l)
		for _, peer := range temp.Peers {
			if k == 0 {
				randomPeer = *peer
			}
			k--
		}
		p.ConnectTo(randomPeer.IP, fmt.Sprint(randomPeer.Port), "interested "+key+"\n")
		it--
	}
}
func (p *Peer) progression(key string, conn net.Conn) {
	i, _ := strconv.Atoi(tools.GetValueFromConfig("Peer", "progress_value"))
	for {
		mutex.Lock()
		previousMessage = "have"
		mutex.Unlock()
		if file, valid := p.Files["self"]; valid {
			length := tools.BitCount(*file.Peers[conn.LocalAddr().String()].BufferMaps[key])
			if length >= i && length%i == 0 {
				WriteReadConnection(conn, p, "have "+key+" "+tools.BufferMapToString(*file.Peers[conn.LocalAddr().String()].BufferMaps[key])+"\n")
			}
		}
	}
}
func (p *Peer) ConnectTo(IP string, Port string, mess ...string) {
	conn, err := net.Dial("tcp", IP+":"+Port)
	errorCheck(err)
	// defer conn.Close()
	if _, valid := p.Comm[conn.RemoteAddr().String()]; !valid {
		p.Comm[conn.RemoteAddr().String()] = conn
		fmt.Println(conn.LocalAddr(), " is connected to ", conn.RemoteAddr())
	}
	if len(mess) == 0 {
		WriteReadConnection(conn, p)
	} else {
		WriteReadConnection(conn, p, mess...)
	}

	// go handleConnection(conn)
	//Handle the response here !
	// fmt.Print(string(buffer[:n]))
}
