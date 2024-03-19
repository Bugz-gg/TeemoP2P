package tools

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var LocalFiles = map[string]File{} // Supposing no collision will happen during the project
var RemoteFiles = map[string]File{}

type File struct {
	Name      string
	Size      int
	PieceSize int
	Key       string
	BufferMap BufferMap
}

type ListData struct {
	Files []File
}

type InterestedData struct {
	Key string
}

type HaveData struct {
	Key       string
	BufferMap BufferMap
}

func ListRegexGen() (ListRegex func() *regexp.Regexp) {
	listPattern := `^list \[(.*)\]$`
	listRegex := regexp.MustCompile(listPattern)
	return func() *regexp.Regexp {
		return listRegex
	}
}

var ListRegex = ListRegexGen()

func InterestedRegexGen() (InterestedRegex func() *regexp.Regexp) {
	interestedPattern := `^interested (?<key>[a-zA-Z0-9]{32})$` // Add optional leech if necessary
	interestedRegex := regexp.MustCompile(interestedPattern)
	return func() *regexp.Regexp {
		return interestedRegex
	}
}

var InterestedRegex = InterestedRegexGen()

func HaveRegexGen() (HaveRegex func() *regexp.Regexp) {
	havePattern := `^have ([a-zA-Z0-9]{32}) ([01]*)$` // Add optional leech if necessary
	haveRegex := regexp.MustCompile(havePattern)
	return func() *regexp.Regexp {
		return haveRegex
	}
}

var HaveRegex = HaveRegexGen()

func GetPiecesRegexGen() (GetPiecesRegex func() *regexp.Regexp) {
	getPiecesPattern := `^getpieces [a-zA-Z0-9]{32} \[[0-9 ]*\]$` // Add optional leech if necessary
	getPiecesRegex := regexp.MustCompile(getPiecesPattern)
	return func() *regexp.Regexp {
		return getPiecesRegex
	}
}

var GetPiecesRegex = GetPiecesRegexGen()

func DataRegexGen() (DataRegex func() *regexp.Regexp) {
	dataPattern := `^getpieces [a-zA-Z0-9]{32} \[(?:[0-9]*:[01]*| )\]$` // Add optional leech if necessary
	dataRegex := regexp.MustCompile(dataPattern)
	return func() *regexp.Regexp {
		return dataRegex
	}
}

var DataRegex = DataRegexGen()

func ListCheck(message string) (bool, ListData) {
	if match := ListRegex().FindStringSubmatch(message); match != nil {
		filesData := strings.Split(match[1], " ")
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
				errors.New("Key error.")
			}
			file := File{Name: filename, Size: size, PieceSize: pieceSize, Key: key, BufferMap: BufferMap{Length: (size-1)/pieceSize/8 + 1, BitSequence: make([]byte, (size-1)/pieceSize/8+1)}}
			listStruct.Files = append(listStruct.Files, file)
			RemoteFiles[key] = file // Update the registered remote files.
		}
		return true, listStruct
	}
	return false, ListData{}
}

func InterestedCheck(message string) (bool, InterestedData) {
	if match := InterestedRegex().FindStringSubmatch(message); match != nil {
		return true, InterestedData{Key: match[1]}
	}
	return false, InterestedData{}
}

func HaveCheck(message string) (bool, HaveData) {
	if match := HaveRegex().FindStringSubmatch(message); match != nil {
		buffer := match[2]
		if len(buffer) == 0 {
			buffer = "0"
		}
		file := RemoteFiles[match[1]]
		file = File{Size: 12, PieceSize: 1, Key: "Uizhsja8hzUizhsja8hzUizhsja8hzsu"} // To be removed.
		if len(buffer) != BufferBitSize(file) {
			return false, HaveData{}
		}
		return true, HaveData{Key: match[1], BufferMap: StringToBufferMap(buffer)}
	}
	return false, HaveData{}
}

func (f *File) GetFile() (string, int, int, string, bool) {
	if f.Name == "" && f.Size == 0 {
		return f.Name, f.Size, f.PieceSize, f.Key, false
	}
	return f.Name, f.Size, f.PieceSize, f.Key, true
}

// Fonction de mise à jour des peers
