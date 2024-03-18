package tools

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	//"os"
	//peer "peerproject/pair"
)

type File struct {
	Name      string
	Size      int
	PieceSize int
	Key       string
	BufferMap BufferMap
}

type BufferMap struct {
	Length      int
	BitSequence []byte
}

type AnnounceData struct {
	Port      int
	Files     []File
	LeechKeys []string
}

type InterestedData struct {
	Key string
}

type HaveData struct {
	Key       string
	BufferMap BufferMap
}

func SearchFile(key string) {

}

func BufferBitSize(file File) int {
	return (file.Size-1)/file.PieceSize + 1
}

func BufferSize(file File) int {
	return (file.Size-1)/file.PieceSize/8 + 1
}

func ByteArrayWrite(array []byte, index int, value int) {
	if value == 1 {
		array[index/8] |= 1 << (7 - (index % 8))
	}
}

func BufferMapWrite(bufferMap BufferMap, index int, value int) {

}

func StringToBufferMap(str string) BufferMap {
	array := make([]byte, (len(str)-1)/8+1)
	for index, char := range str {
		if char == '1' {
			ByteArrayWrite(array, index, 1)
		}
	}
	return BufferMap{Length: len(str), BitSequence: array}
}

func BufferMapToString(bufferMap BufferMap) {

}

func AnnounceRegexGen() (AnounceRegex func() *regexp.Regexp) {
	announcePattern := `^announce\s+listen\s+(\d+)\s+seed\s+\[(.*)\]$` // Add optional leech if necessary
	announceRegex := regexp.MustCompile(announcePattern)
	return func() *regexp.Regexp {
		return announceRegex
	}
}

var AnnounceRegex = AnnounceRegexGen()

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

func (f *File) GetFile() (string, int, int, string, bool) {
	if f.Name == "" && f.Size == 0 {
		return f.Name, f.Size, f.PieceSize, f.Key, false
	}
	return f.Name, f.Size, f.PieceSize, f.Key, true
}

func RegexInit() (*regexp.Regexp, *regexp.Regexp, *regexp.Regexp) {
	announcePattern := `^announce\s+listen\s+(\d+)\s+seed\s+\[(.*)\]$` // Add optional leech if necessary
	announceRegex := regexp.MustCompile(announcePattern)

	lookPattern := `^look\s+\[[\w+(?:<|<=|!=|=|>|>=)\".*\"\s+]*\]$`
	lookRegex := regexp.MustCompile(lookPattern)

	getfilePattern := `^getfile\s+[a-z0-9]{32}$`
	getfileRegex := regexp.MustCompile(getfilePattern)

	return announceRegex, lookRegex, getfileRegex
}

func AnnounceCheck(message string) (bool, AnnounceData) {
	if match := AnnounceRegex().FindStringSubmatch(message); match != nil {
		//captured := len(match)-1
		port, _ := strconv.Atoi(match[1])
		filesData := strings.Split(match[2], " ")
		if len(filesData)%4 != 0 {
			fmt.Println("Invalid received message.")
			return false, AnnounceData{}
		}
		announceStruct := AnnounceData{Port: port}
		nbFiles := len(filesData) / 4
		for i := 0; i < nbFiles; i++ {
			filename := filesData[i*4]
			size, err := strconv.Atoi(filesData[i*4+1])
			pieceSize, err2 := strconv.Atoi(filesData[i*4+2])

			if err != nil || err2 != nil {
				fmt.Println("Invalid conversion to int (size or piece size).", err, err2)
				return false, AnnounceData{}
			}
			key := filesData[i*4+3]
			if len(key) != 32 {
				errors.New("Key error.")
			}
			announceStruct.Files = append(announceStruct.Files, File{Name: filename, Size: size, PieceSize: pieceSize, Key: key, BufferMap: BufferMap{Length: (size-1)/pieceSize/8 + 1, BitSequence: make([]byte, (size-1)/pieceSize/8+1)}})
		}
		fmt.Println(announceStruct.Files)
	}
	return false, AnnounceData{}
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
		file := File{Size: 12, PieceSize: 1, Key: "Uizhsja8hzUizhsja8hzUizhsja8hzsu"}
		if len(buffer) != BufferBitSize(file) {
			return false, HaveData{}
		}
		return true, HaveData{Key: match[1], BufferMap: StringToBufferMap(buffer)}
	}
	return false, HaveData{}
}

// Fonction de mise à jour des peers
