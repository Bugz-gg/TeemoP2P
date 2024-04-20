package tools

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// LocalFiles A map to store local files' data.
var LocalFiles *map[string]*File // Supposing no collision will happen during the project

// RemoteFiles A map to store remote files' data.
var RemoteFiles = map[string]*File{} // Supposing no collision will happen during the project

// AllPeers A map to store all connected peers.
var AllPeers = map[string]*Peer{}

//var RemotePeerFiles map[string]PeersData

// AddFile adds a file to a map. (LocalFiles and RemoteFiles in this project.)
func AddFile(fileMap *map[string]*File, file *File) {
	(*fileMap)[file.Key] = file
}

// RemoveFile removes a file from a map.  (LocalFiles and RemoteFiles in this project.)
func RemoveFile(fileMap *map[string]*File, file File) {
	delete(*fileMap, file.Key)
}

// Map is a function to apply a function on all elements of an array.
func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}

// ListRegexGen provides the function that returns the compiled regex expression for the `list` message.
func ListRegexGen() (ListRegex func() *regexp.Regexp) {
	listPattern := `^list \[(.*)\]$`
	listRegex := regexp.MustCompile(listPattern)
	return func() *regexp.Regexp {
		return listRegex
	}
}

// ListRegex is the function provided by ListRegexGen.
var ListRegex = ListRegexGen()

// InterestedRegexGen provides the function that returns the compiled regex expression for the `interested` message.
func InterestedRegexGen() (InterestedRegex func() *regexp.Regexp) {
	interestedPattern := `^interested ([a-zA-Z0-9]{32})$` // Add optional leech if necessary
	interestedRegex := regexp.MustCompile(interestedPattern)
	return func() *regexp.Regexp {
		return interestedRegex
	}
}

// InterestedRegex is the function provided by InterestedRegexGen.
var InterestedRegex = InterestedRegexGen()

// HaveRegexGen provides the function that returns the compiled regex expression for the `have` message.
func HaveRegexGen() (HaveRegex func() *regexp.Regexp) {
	havePattern := `^have ([a-zA-Z0-9]{32}) ([01]*)$` // Add optional leech if necessary
	haveRegex := regexp.MustCompile(havePattern)
	return func() *regexp.Regexp {
		return haveRegex
	}
}

// HaveRegex is the function provided by HaveRegexGen.
var HaveRegex = HaveRegexGen()

// GetPiecesRegexGen provides the function that returns the compiled regex expression for the `getpieces` message.
func GetPiecesRegexGen() (GetPiecesRegex func() *regexp.Regexp) {
	getPiecesPattern := `^getpieces ([a-zA-Z0-9]{32}) \[([0-9 ]*)\]$` // Add optional leech if necessary
	getPiecesRegex := regexp.MustCompile(getPiecesPattern)
	return func() *regexp.Regexp {
		return getPiecesRegex
	}
}

// GetPiecesRegex is the function provided by GetPiecesRegexGen.
var GetPiecesRegex = GetPiecesRegexGen()

// DataRegexGen provides the function that returns the compiled regex expression for the `data` message.
func DataRegexGen() (DataRegex func() *regexp.Regexp) { // To be tested
	dataPattern := `^data ([a-zA-Z0-9]{32}) \[((?:[0-9]*:[01]*| )*)\]$` // Add optional leech if necessary
	dataRegex := regexp.MustCompile(dataPattern)
	return func() *regexp.Regexp {
		return dataRegex
	}
}

// DataRegex is the function provided by DataRegexGen.
var DataRegex = DataRegexGen()

// PeersRegexGen provides the function that returns the compiled regex expression for the `peers` message.
func PeersRegexGen() (PeersRegex func() *regexp.Regexp) { // IPv4
	peersPattern := `^peers ([a-zA-Z0-9]{32}) \[((?:[0-9]+.[0-9]+.[0-9]+.[0-9]+:[0-9]+| )*)\]$`
	peersRegex := regexp.MustCompile(peersPattern)
	return func() *regexp.Regexp {
		return peersRegex
	}
}

// PeersRegex is the function provided by PeersRegexGen.
var PeersRegex = PeersRegexGen()

// ListCheck checks the format of a `list` message. The boolean tells whether the format is valid or not. The returned struct's validity depends on the boolean.
func ListCheck(message string) (bool, ListData) {
	if match := ListRegex().FindStringSubmatch(message); match != nil {
		var filesData []string
		if match[1] != "" {
			filesData = strings.Split(match[1], " ")
		}
		if len(filesData)%4 != 0 {
			fmt.Println("Invalid received message.")
			return false, ListData{}
		}
		listStruct := ListData{}
		nbFiles := len(filesData) / 4
		for i := 0; i < nbFiles; i++ {
			filename := filesData[i*4]
			size, err := strconv.Atoi(filesData[i*4+1])
			pieceSize, err2 := strconv.Atoi(filesData[i*4+2])

			if err != nil || err2 != nil {
				fmt.Println("Invalid conversion to int (size or piece size).", err, err2)
				return false, ListData{}
			}
			key := filesData[i*4+3]
			if len(key) != 32 {
				fmt.Println("Invalid key.", err, err2)
				return false, ListData{}
			}
			if _, valid := RemoteFiles[key]; !valid { // If not registered as a RemoteFile.
				RemoteFiles[key] = &File{Name: filename, Size: size, PieceSize: pieceSize, Key: key, Peers: make(map[string]*Peer)} // Update the registered remote files.
			}
			file := RemoteFiles[key]
			// file := File{Name: filename, Size: size, PieceSize: pieceSize, Key: key} //, BufferMap: BufferMap{Length: size / pieceSize, BitSequence: make([]byte, (size-1)/pieceSize/8+1)}}
			listStruct.Files = append(listStruct.Files, *file)
			//RemoteFiles[key] = &file
		}
		return true, listStruct
	}
	return false, ListData{}
}

// InterestedCheck checks the format of a `interested` message. The boolean tells whether the format is valid or not. The returned struct's validity depends on the boolean.
func InterestedCheck(message string) (bool, InterestedData) {
	if match := InterestedRegex().FindStringSubmatch(message); match != nil {
		if _, valid := (*LocalFiles)[match[1]]; !valid {
			fmt.Println("No such file locally.")
			return false, InterestedData{}
		}
		return true, InterestedData{Key: match[1]}
	}
	return false, InterestedData{}
}

// HaveCheck checks the format of a `have` message. The boolean tells whether the format is valid or not. The returned struct's validity depends on the boolean.
func HaveCheck(message string) (bool, HaveData) {
	if match := HaveRegex().FindStringSubmatch(message); match != nil {
		buffer := match[2]
		if len(buffer) == 0 {
			buffer = "0"
		}
		//file := RemoteFiles[match[1]]

		file := RemoteFiles[match[1]] // Change from LocalFiles to RemoteFiles. Maybe add check Local ?
		//file = &File{Size: 12, PieceSize: 1, Key: "Uizhsja8hzUizhsja8hzUizhsja8hzsu"} // To be removed.
		if len(buffer) != BufferBitSize(*file) {
			return false, HaveData{}
		}
		return true, HaveData{Key: match[1], BufferMap: StringToBufferMap(buffer)}
	}
	return false, HaveData{}
}

// GetPiecesCheck checks the format of a `getpieces` message. The boolean tells whether the format is valid or not. The returned struct's validity depends on the boolean.
func GetPiecesCheck(message string) (bool, GetPiecesData) {
	if match := GetPiecesRegex().FindStringSubmatch(message); match != nil {
		buffer := match[2]
		if len(buffer) == 0 {
			fmt.Println("No piece requested.")
			return false, GetPiecesData{}
		}
		if _, valid := (*LocalFiles)[match[1]]; !valid {
			fmt.Println("No such file locally.")
			return false, GetPiecesData{}
		}
		wantedPieces := Map(strings.Split(match[2], " "), func(item string) int { i, _ := strconv.Atoi(item); return i })
		file := (*LocalFiles)[match[1]]
		var pieces []int
		for _, i := range wantedPieces {
			if i < BufferBitSize(*file) && ByteArrayCheck((*LocalFiles)[match[1]].Peers["self"].BufferMaps[match[1]].BitSequence, i) {
				pieces = append(pieces, i)
			} else {
				fmt.Println("Invalid pieces' numbers :", i)
			}
		}
		return true, GetPiecesData{Key: match[1], Pieces: pieces}
	}
	return false, GetPiecesData{}
}

// DataCheck checks the format of a `data` message. The boolean tells whether the format is valid or not. The returned struct's validity depends on the boolean.
func DataCheck(message string) (bool, DataData) {
	if match := DataRegex().FindStringSubmatch(message); match != nil {
		buffer := match[2]
		if len(buffer) == 0 {
			fmt.Println("No piece given.")
			return false, DataData{}
		}
		if _, valid := (*LocalFiles)[match[1]]; !valid { // If we don't have any piece of the requested file yet.
			rFile := RemoteFiles[match[1]]
			(*LocalFiles)[match[1]] = &File{Name: rFile.Name, Size: rFile.Size, PieceSize: rFile.PieceSize, Key: match[1]}
			//InitBufferMap((*LocalFiles)[match[1]])
		}
		piecesdata := strings.Split(match[2], " ")

		file := (*LocalFiles)[match[1]]
		piecesize := file.PieceSize
		pieces := make([]Piece, len(piecesdata))
		for i, data := range piecesdata {
			piece := strings.Split(data, ":")
			index, _ := strconv.Atoi(piece[0])
			if index < 0 || index >= BufferBitSize(*file) { //file.BufferMapLength {
				fmt.Printf("Out or range index received. (%d)\n", index)
				return false, DataData{}
			}
			// Make initialize self with id ?
			if ByteArrayCheck((*LocalFiles)[match[1]].Peers["self"].BufferMaps[match[1]].BitSequence, index) {
				fmt.Printf("Already have the piece at index %d.\n", index)
				//return false, DataData{}
				continue
			}
			if len(piece[1]) != piecesize*8 {
				fmt.Println("Wrong piece size received.")
				return false, DataData{}
			}
			pieces[i].Index = index
			pieces[i].Data = StringToData(piece[1])

			// Check integrity of file if all pieces have been downloaded ?

		}
		return true, DataData{Key: match[1], Pieces: pieces}
	}
	return false, DataData{}
}

// PeersCheck checks the format of a `peers` message. The boolean tells whether the format is valid or not. The returned struct's validity depends on the boolean.
func PeersCheck(message string) bool {
	if match := PeersRegex().FindStringSubmatch(message); match != nil {
		buffer := match[2]
		if len(buffer) == 0 {
			fmt.Println("No peer given.")
			return false
		}
		if _, valid := RemoteFiles[match[1]]; !valid { // If the file is not registered yet (we don't have any information about it).
			fmt.Printf("No data about file %s.", match[1])
			return false
		}

		peersdata := strings.Split(match[2], " ")

		peers := RemoteFiles[match[1]].Peers
		if peers == nil {
			peers = make(map[string]*Peer)
		}

		// Check if peer already registered
		for _, data := range peersdata {
			info := strings.Split(data, ":")
			port := info[1]
			peerId := fmt.Sprintf("%s:%s", info[0], port)
			if _, valid := AllPeers[peerId]; !valid { // If it is the first time learning about a peer.
				AllPeers[peerId] = &Peer{IP: info[0], Port: port}
			}
			if _, valid := RemoteFiles[match[1]].Peers[peerId]; !valid { // Add peer to owners of the remote file if not already in.
				RemoteFiles[match[1]].Peers[peerId] = AllPeers[peerId]
			}

			peers[peerId] = AllPeers[peerId] // Add peer to list of peers having the file.
		}
		return true
	}
	return false
}

func (f *File) GetFile() (string, int, int, string, bool) {
	if f.Name == "" && f.Size == 0 {
		return f.Name, f.Size, f.PieceSize, f.Key, false
	}
	return f.Name, f.Size, f.PieceSize, f.Key, true
}

func (f *File) GetFileUpdate() (string, []byte, bool) {
	if buff, valid := f.Peers["self"].BufferMaps[f.Key]; valid {
		if f.Name == "" && f.Size == 0 {
			return f.Key, buff.BitSequence, false
		}
		return f.Key, buff.BitSequence, true
	} else {
		return f.Key, make([]byte, 0), true
	}
}

// Fonction de mise à jour des peers
