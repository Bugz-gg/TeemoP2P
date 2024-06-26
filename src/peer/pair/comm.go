package peer_package

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"peerproject/tools"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var mutex sync.Mutex
var previousMessage string // To know if it was an interested or a have.<ScrollWheelDown>
var rare = false
var dl = false

var DlDone chan struct{}
var DlDoneComm chan struct{}
var ResponsesRemainingUpdated chan struct{}
var ResponsesRemaining atomic.Int64
var DlFile map[string]*os.File // Currently downloading file and manifest
var dlKey string               // And key

func WaitFor(sig *chan struct{}, timeoutRemaining bool, timeout bool, timeoutTime time.Duration) bool {
	var timer <-chan time.Time
	if timeout {
		timer = time.After(timeoutTime)
	} else if timeoutRemaining {
		<-ResponsesRemainingUpdated
		timePerMessage, _ := strconv.ParseInt(tools.GetValueFromConfig("Peer", "response_timeout"), 10, 64)
		timer = time.After(time.Duration(ResponsesRemaining.Load()*timePerMessage) * time.Second)
	}

	select {
	case <-*sig:
		return true
	case <-timer:
		return false
	}
}

func ChannSignal(sig *chan struct{}) {
	*sig <- struct{}{}
	*sig = make(chan struct{})
}

func (p *Peer) HelloTrack(t Peer) {
	timeout, _ := strconv.Atoi(tools.GetValueFromConfig("Peer", "timeout"))
	message := "announce listen " + p.Port + " seed ["
	leechFiles := []*tools.File{}
	for _, file := range p.Files {
		if !file.Complete {
			leechFiles = append(leechFiles, file)
			continue
		}
		name, size, pieceSize, key, isEmpty := file.GetFile()
		if isEmpty {
			message += fmt.Sprintf(`%s %d %d %s `, name, size, pieceSize, key)
		} else {
			break
		}
	}
	message = strings.TrimSuffix(message, " ") + "]"
	if len(leechFiles) > 0 {
		message += " leech ["
		for _, file := range leechFiles {
			message += fmt.Sprintf(`%s `, file.Key)
		}
		message = strings.TrimSuffix(message, " ") + "]"
	}
	message += "\n"
	conn, err := net.Dial("tcp", t.IP+":"+t.Port)
	errorCheck(err)
	// defer conn.Close()
	_, err = conn.Write([]byte(message))
	errorCheck(err)
	buffer := make([]byte, 1024) // 3 would be sufficient.
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

func (p *Peer) sendupdate() {
	var seed string
	var leech string
	timeout, _ := strconv.Atoi(tools.GetValueFromConfig("Peer", "timeout"))
	updateValue, _ := strconv.Atoi(tools.GetValueFromConfig("Peer", "update_time"))
	for {
		seed = ""
		leech = ""
		time.Sleep(time.Duration(float64(updateValue) * math.Pow(10, 9)))
		message := "update seed ["
		for _, file := range p.Files {
			key, BitSequence, notEmpty := file.GetFileUpdate()
			if notEmpty {
				temp := 0
				for k := range len(BitSequence) {
					if !tools.ByteArrayCheck(BitSequence, uint64(k)) {
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
		conn := p.Comm["tracker"]
		_, err := conn.Write([]byte(message))
		errorCheck(err)
		buffer := make([]byte, 1024)
		err = conn.SetReadDeadline(time.Now().Add(time.Duration(float64(timeout) * math.Pow(10, 9))))
		n, err := conn.Read(buffer)
		err = conn.SetReadDeadline(time.Time{})
		if err != nil {
			fmt.Print(string(buffer[:n]))
		}
	}
}

func (p *Peer) rarepiece() {
	conn := p.Comm["tracker"]
	t, _ := strconv.Atoi(tools.GetValueFromConfig("Peer", "time_dl_rare_piece"))
	var byteArray []int
	var connArray map[int][]net.Conn
	for {
		time.Sleep(time.Duration(math.Pow(10, 9) * float64(t)))
		WriteReadConnection(conn, p, "look []\n")
		for key := range tools.RemoteFiles {
			// fmt.Println("Im in the boucle :)")
			if _, valid := p.Files[key]; !valid || !tools.ArrayCheck(*p.Files[key].Peers["self"].BufferMaps[key]) {
				// fmt.Println("rare piece key :", key)
				byteArray = make([]int, tools.BufferBitSize(*tools.RemoteFiles[key]))
				connArray = make(map[int][]net.Conn, tools.BufferBitSize(*tools.RemoteFiles[key]))
				// fmt.Println("byte array: ", byteArray, connArray)
				k := max(int(tools.BufferBitSize(*tools.RemoteFiles[key])/4), 1)
				// fmt.Println(k, tools.RemoteFiles[key].PieceSize, tools.RemoteFiles[key].Size)
				rareArray := make([]int, k)
				mutex.Lock()
				rare = true
				mutex.Unlock()
				WriteReadConnection(conn, p, "getfile "+key+"\n")
				for connRem := range tools.RemoteFiles[key].Peers {
					// TODO : Faire un if pour regarder si on a deja fait un interested en regardant p.BufferMaps
					if rconn, valid := p.Comm[connRem]; !valid {
						p.ConnectTo(tools.RemoteFiles[key].Peers[connRem].IP, tools.RemoteFiles[key].Peers[connRem].Port, "interested "+key+"\n")
					} else {
						WriteReadConnection(rconn, p, "interested "+key+"\n")
					}
					for i := range tools.RemoteFiles[key].Peers[connRem].BufferMaps[key].Length {
						if byteArray[i] != int(math.Inf(1)) {
							if tools.ByteArrayCheck(p.Files[key].Peers["self"].BufferMaps[key].BitSequence, i) {
								byteArray[i] = int(math.Inf(1))
							} else if tools.ByteArrayCheck(tools.RemoteFiles[key].Peers[connRem].BufferMaps[key].BitSequence, i) {
								connArray[int(i)] = append(connArray[int(i)], p.Comm[connRem])
								byteArray[i] += 1
							}
						}
					}

					// fmt.Print(tools.RemoteFiles[key].Peers[connRem].BufferMaps[key].BitSequence, byteArray, connArray, tools.RemoteFiles[key].Peers[connRem].BufferMaps[key].Length)
				}
				minIndex := 0
				for i := 0; i < k; i++ {
					minIndex = i
					for x := i + 1; x < len(byteArray); x++ {
						if (byteArray[minIndex] < 0) || (byteArray[x] > 0 && byteArray[x] < byteArray[minIndex]) {
							minIndex = x
						}
					}
					// fmt.Println(i, minIndex)
					byteArray[i], byteArray[minIndex] = byteArray[minIndex], byteArray[i]
					rareArray[i] = minIndex
				}
				for _, index := range rareArray { // TODO : Can improve in case its only one peer to send it only once
					if index < 0 {
						continue
					}
					// fmt.Println(connArray, connArray[index], len(connArray[index]), rareArray, byteArray)
					WriteReadConnection(connArray[index][rand.Intn(len(connArray[index]))], p, "getpieces "+key+" ["+strconv.Itoa(index)+"]\n") // Beter with sprintf maybe
				}
				break
			}
		}
	}

}

func (p *Peer) Downloading(key string) {
	buffSize, _ := strconv.ParseUint(tools.GetValueFromConfig("Peer", "max_buff_size"), 10, 64)
	if tools.RemoteFiles[key].PieceSize > buffSize+100 { // Le 50 sert au format du message `data`
		fmt.Printf("Impossible de télécharger \033[0;35m%s\033[0m. La taille des pièces est trop grosse.", p.Files[key].Name) // Normalement impossible d'arriver ici.
		tools.WriteLog("Impossible de télécharger %s. La taille des pièces est trop grosse.", p.Files[key].Name)
		return
	}
	mutex.Lock()
	dl = true
	mutex.Unlock()
	indexByConn := map[uint64][]net.Conn{} // All the indexes we don't have yet and the array of connections that have them.
	connAsk := map[net.Conn][]string{}     // The indexes we'll ask from each connection.
	var dontHave []uint64                  // To sort the wanted indexes by ascending order of number of peers having it.
	if _, valid := p.Files[key]; valid {
		for index := range p.Files[key].Peers["self"].BufferMaps[key].Length {
			if !tools.ByteArrayCheck(p.Files[key].Peers["self"].BufferMaps[key].BitSequence, index) { // Check which pieces are missing.
				indexByConn[index] = []net.Conn{}
				dontHave = append(dontHave, index)
			}
		}
	} else {
		for i := range tools.BufferBitSize(*tools.RemoteFiles[key]) {
			dontHave = append(dontHave, i)
		}
	}
	WriteReadConnection(p.Comm["tracker"], p, "getfile "+key+"\n") // Update the peers having the file.
	for connRem := range tools.RemoteFiles[key].Peers {
		if tools.RemoteFiles[key].Peers[connRem].IP == p.IP && tools.RemoteFiles[key].Peers[connRem].Port == p.Port { // No need to ask pieces from ourselves
			continue
		}
		if rconn, valid := p.Comm[connRem]; !valid {
			p.ConnectTo(tools.RemoteFiles[key].Peers[connRem].IP, tools.RemoteFiles[key].Peers[connRem].Port, "interested "+key+"\n")
		} else {
			WriteReadConnection(rconn, p, "interested "+key+"\n")
		}
		for index := range dontHave { // Get the list peers having a certain missing piece.
			uintindex := uint64(index)
			if tools.ByteArrayCheck(tools.RemoteFiles[key].Peers[connRem].BufferMaps[key].BitSequence, uintindex) {
				indexByConn[uintindex] = append(indexByConn[uintindex], p.Comm[connRem])
			}
		}
	}

	sort.SliceStable(dontHave, func(i, j int) bool { // Sort by ascending order of number of peers.
		return len(indexByConn[dontHave[i]]) < len(indexByConn[dontHave[j]])
	})

	for index := range dontHave { // Try to share the work between peers for each wanted piece.
		uintindex := uint64(index)
		if len(indexByConn[uintindex]) == 0 { // No peer has the piece.
			continue
		}
		conns := indexByConn[uintindex] // List of peers having the piece.
		requested := conns[0]
		for _, conn := range conns[1:] {
			if len(connAsk[conn]) > len(connAsk[requested]) {
				break
			}
			if len(connAsk[conn]) < len(connAsk[requested]) {
				requested = conn
				break
			}
		}
		connAsk[requested] = append(connAsk[requested], strconv.Itoa(index))
	}

	var nbPieces = int64((buffSize - 100) / p.Files[key].PieceSize)
	for _, indexes := range connAsk {
		pieces := int64(len(indexes))
		ResponsesRemaining.Add((pieces + nbPieces - 1) / nbPieces)
	}
	ChannSignal(&ResponsesRemainingUpdated)
	dlKey = key
	path := tools.GetValueFromConfig("Peer", "path")

	filename := tools.RemoteFiles[key].Name
	os.MkdirAll(filepath.Join("./", path, "/", filename), os.FileMode(0777))
	DlFile["file"], _ = os.OpenFile(filepath.Join("./", path, filename+"/file"), os.O_CREATE|os.O_RDWR, os.FileMode(0777))
	DlFile["manifest"], _ = os.OpenFile(filepath.Join("./", path, filename+"/manifest"), os.O_CREATE|os.O_RDWR, os.FileMode(0777))

	for conn, indexes := range connAsk {
		go func() {
			var i int64 = 0
			for ; i < int64(len(indexes)); i += nbPieces {
				WriteReadConnection(conn, p, "getpieces "+key+" ["+strings.Join(indexes[i:min(len(indexes), int(i+nbPieces))], " ")+"]\n")
			}
		}()
	}
	bufferSize := float64(tools.BufferBitSize(*p.Files[key]))
	bufferMap := *p.Files[key].Peers["self"].BufferMaps[key]
	go func() { // Maybe just maybe can improve this but it's ok for now or I'm sure now
		currPercent := int(float64(tools.BitCount(bufferMap)) * 100 / bufferSize)
		str := "▌"
		fmt.Print(strings.Repeat(str, currPercent/10))
		for !tools.ArrayCheck(*p.Files[key].Peers["self"].BufferMaps[key]) {
			newPercent := int(float64(tools.BitCount(bufferMap)) * 100 / bufferSize)
			fmt.Print(strings.Repeat(str, (int(newPercent)/10)-(currPercent/10)))
			currPercent = newPercent
		}
		newPercent := int(float64(tools.BitCount(bufferMap)) * 100 / bufferSize)
		fmt.Println(strings.Repeat(str, (int(newPercent)/10)-(currPercent/10)))
	}()
	WaitFor(&DlDone, false, true, time.Second*time.Duration(ResponsesRemaining.Load()))
	ChannSignal(&DlDoneComm)
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
	if message == "exit\n" {
		os.Exit(1)
	}
	conn.Write([]byte(message))
	buffSize, _ := strconv.Atoi(tools.GetValueFromConfig("Peer", "max_buff_size"))
	buffer := make([]byte, buffSize)
	var eom string
	var n int
	var nerr error = nil
	var fd int
	for len(eom) == 0 || eom[len(eom)-1] != '\n' || nerr != nil {
		conn.SetReadDeadline(time.Now().Add(time.Duration(float64(20) * math.Pow(10, 9))))
		fd, nerr = conn.Read(buffer)
		n += fd
		// errorCheck(nerr)
		eom += string(buffer[:fd])
	}
	conn.SetReadDeadline(time.Time{})
	if n > 0 {
		mess := eom
		//mess = strings.TrimSuffix(mess, "\n")
		input := strings.Split(mess, " ")[0]
		// fmt.Printf("[\u001B[0;33m%s\u001B[39m]:%s\n", conn.RemoteAddr().String(), mess)
		//tools.WriteLog("%s:%s\n", conn.RemoteAddr().String(), mess)
		switch input {
		case "data":
			valid, data := tools.DataCheck(mess)
			if valid {
				var fdf *os.File
				var fdc *os.File
				downloading := false
				file := p.Files[data.Key]
				path := tools.GetValueFromConfig("Peer", "path")

				if data.Key == dlKey && DlFile["file"] != nil {
					fdf = DlFile["file"]
					fdc = DlFile["manifest"]
					downloading = true
				} else {
					os.MkdirAll(filepath.Join("./", path, "/", p.Files[data.Key].Name), os.FileMode(0777))

					fdf, _ = os.OpenFile(filepath.Join("./", path, file.Name+"/file"), os.O_CREATE|os.O_RDWR, os.FileMode(0777))
					//errorCheck(err)
					fdc, _ = os.OpenFile(filepath.Join("./", path, file.Name+"/manifest"), os.O_CREATE|os.O_RDWR, os.FileMode(0777))
					//errorCheck(err)
				}
				bufferMap := file.Peers["self"].BufferMaps[data.Key]
				fdc.Seek(0, 0)
				_, _ = fdc.WriteString(file.Name + "\n")
				_, _ = fdc.WriteString(strconv.FormatUint(file.Size, 10) + "\n")
				_, _ = fdc.WriteString(strconv.FormatUint(file.PieceSize, 10) + "\n")
				_, err := fdc.WriteString(file.Key + "\n")
				bufferMapOffset, _ := fdc.Seek(0, io.SeekCurrent)
				errorCheck(err)
				for i := 0; len(data.Pieces) > i; i++ {
					_, err := fdf.Seek(int64(data.Pieces[i].Index*file.PieceSize), 0)
					errorCheck(err)
					var n int
					if file.Size%file.PieceSize != 0 && data.Pieces[i].Index+1 == tools.BufferBitSize(*file) {
						n, err = fdf.Write(data.Pieces[i].Data.String[:file.Size%file.PieceSize])
						errorCheck(err)

					} else {
						n, err = fdf.Write(data.Pieces[i].Data.String)
						errorCheck(err)
					}
					if n <= 0 {
						fmt.Println("File has not being written. :(", err) //, data.Pieces[i].Data.String)
						continue
					}
					tools.BufferMapWrite(bufferMap, data.Pieces[i].Index)
				}
				//_, err = fdc.WriteString(tools.BufferMapToString(*file.Peers["self"].BufferMaps[data.Key]) + "\n")

				for _, piece := range data.Pieces {
					fdc.Seek(bufferMapOffset+int64(piece.Index), io.SeekStart)
					_, _ = fdc.WriteString("1")
				}
				fdc.Seek(bufferMapOffset+int64(tools.BufferBitSize(*file)), io.SeekStart)
				_, _ = fdc.WriteString("\n")

				if !downloading {
					fdf.Close()
					fdc.Close()
				}

				errorCheck(err)
				if bufferMap.Count >= tools.BufferBitSize(*file) && tools.GetMD5Hash(filepath.Join(path, file.Name+"/file")) == data.Key {
					if _, valid := ServerOpenedFiles[data.Key]; valid {
						ServerOpenedFiles[data.Key].File.Close()
						delete(ServerOpenedFiles, data.Key)
					}
					fdf.Close()
					fdc.Close()
					err = os.Rename(filepath.Join("./", path, file.Name), filepath.Join(path, "todelete"+file.Key))
					errorCheck(err)
					err = os.Rename(filepath.Join("./", path, "todelete"+file.Key, "/file"), filepath.Join(path, file.Name))
					errorCheck(err)

					err = os.RemoveAll(filepath.Join(path, "todelete"+file.Key))
					errorCheck(err)
					tools.AddFile(tools.LocalFiles, file)
					tools.RemoveFile(&tools.RemoteFiles, file)
					fmt.Printf("\u001B[92m%s [%s] Download complete\u001B[39m\n", file.Name, file.Key)
				}

				resp := ResponsesRemaining.Add(-1)

				if resp < 1 {
					go ChannSignal(&DlDone)
				}
			} else {
				//fmt.Printf("\u001B[92mInvalid data response.\u001B[39m: %s\n", mess)
				tools.WriteLog("Invalid data response.\n")
			}
		case "have":
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
					// if _, valid := fil.Peers["self"].BufferMaps[data.Key]; !valid {
					// 	p.Files[data.Key].Peers["self"].BufferMaps = make(map[string]*tools.BufferMap)
					// }
					fil.Peers[conn.LocalAddr().String()] = fil.Peers["self"]
					p.Files[data.Key] = &fil
					// fmt.Println(p.Files[data.Key].Peers[conn.LocalAddr().String()].BufferMaps[data.Key])
				}
				if previousMessage == "interested" && !rare && !dl {
					go p.progression(data.Key, conn)
					time.Sleep(2)
				} else if rare {
					time.Sleep(2)
					mutex.Lock()
					rare = false
					mutex.Unlock()
				} else if dl {
					mutex.Lock()
					dl = false
					mutex.Unlock()
				}

			} else {
				fmt.Printf("\u001B[92mInvalid have response.\u001B[39m: %s\n\n", mess)
				tools.WriteLog("Invalid have response.\n")
			}
		case "OK":

		case "list":
			valid, _ := tools.ListCheck(mess)
			if !valid {
				fmt.Printf("\u001B[92mInvalid list response.\u001B[39m: %s\n", mess)
				tools.WriteLog("Invalid list response.\n")
			}
		case "peers":
			valid := tools.PeersCheck(mess, p.IP+":"+p.Port)
			if valid {
				key := strings.Split(mess, " ")[1]
				mutex.Lock()
				previousMessage = "interested"
				mutex.Unlock()
				if !rare && !dl {
					p.interested(key)
				}
			} else {
				fmt.Printf("\u001B[92mInvalid peers response.\u001B[39m%s\n", mess)
				tools.WriteLog("Invalid peers response.\n")
			}
		case "\u001B[92mInvalid command you have no tries remaining, connection is closed...\u001B[39m":
			p.Close(conn.RemoteAddr().String())
		default:
			// panic("valeur par default et pas parmi la liste")
			// je dois prendre en compte que si je n'ai plus d'essai de fermer de mon cote la conn.
		}
	}
}

func (p *Peer) interested(key string) {
	temp := tools.RemoteFiles[key]
	l := len(temp.Peers)
	it := l
	for max, _ := strconv.Atoi(tools.GetValueFromConfig("Peer", "max_peers_to_connect")); max != 0 && it > 0; max-- {
		var randomPeer tools.Peer
		k := rand.Intn(l)
		for _, peer := range temp.Peers {
			if k == 0 {
				if peer.IP == p.IP && peer.Port == p.Port {
					continue
				}
				randomPeer = *peer
			}
			k--
		}
		if _, valid := p.Comm[randomPeer.IP+":"+randomPeer.Port]; !valid {
			p.ConnectTo(randomPeer.IP, fmt.Sprint(randomPeer.Port), "interested "+key+"\n")
		} else {
			WriteReadConnection(p.Comm[randomPeer.IP+":"+randomPeer.Port], p, "interested "+key+"\n")

		}
		it--
	}
}
func (p *Peer) progression(key string, conn net.Conn) {
	i, _ := strconv.Atoi(tools.GetValueFromConfig("Peer", "progress_value"))
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
func (p *Peer) ConnectTo(IP string, Port string, mess ...string) {
	fmt.Println(p, IP, Port)
	if IP+":"+Port == p.IP+":"+p.Port {
		fmt.Println("\033[0;31mCan't connect to yourself !\033[0m")
		return
	}
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
