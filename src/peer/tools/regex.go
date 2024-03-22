package tools

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var LocalFiles = map[string]*File{}  // Supposing no collision will happen during the project
var RemoteFiles = map[string]*File{} // Switch to map[string]*File{} ?

func AddFile(fileMap map[string]*File, file *File) {
	fileMap[file.Key] = file
}

func RemoveFile(fileMap map[string]*File, file File) {
	delete(fileMap, file.Key)
}

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
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
	getPiecesPattern := `^getpieces ([a-zA-Z0-9]{32}) \[([0-9 ]*)\]$` // Add optional leech if necessary
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
				fmt.Println("Invalid key.", err, err2)
				return false, ListData{}
			}
			file := File{Name: filename, Size: size, PieceSize: pieceSize, Key: key, BufferMap: BufferMap{Length: (size-1)/pieceSize/8 + 1, BitSequence: make([]byte, (size-1)/pieceSize/8+1)}}
			listStruct.Files = append(listStruct.Files, file)
			RemoteFiles[key] = &file // Update the registered remote files.
		}
		return true, listStruct
	}
	return false, ListData{}
}

func InterestedCheck(message string) (bool, InterestedData) {
	if match := InterestedRegex().FindStringSubmatch(message); match != nil {
		if _, valid := LocalFiles[match[1]]; !valid {
			fmt.Println("No such file locally.")
			return false, InterestedData{}
		}
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
		file = &File{Size: 12, PieceSize: 1, Key: "Uizhsja8hzUizhsja8hzUizhsja8hzsu"} // To be removed.
		if len(buffer) != BufferBitSize(*file) {
			return false, HaveData{}
		}
		return true, HaveData{Key: match[1], BufferMap: StringToBufferMap(buffer)}
	}
	return false, HaveData{}
}

func GetPiecesCheck(message string) (bool, GetPiecesData) {
	if match := GetPiecesRegex().FindStringSubmatch(message); match != nil {
		buffer := match[2]
		if len(buffer) == 0 {
			fmt.Println("No piece requested.")
			return false, GetPiecesData{}
		}
		if _, valid := LocalFiles[match[1]]; !valid {
			fmt.Println("No such file locally.")
			return false, GetPiecesData{}
		}
		pieces := Map(strings.Split(match[2], " "), func(item string) int { i, _ := strconv.Atoi(item); return i })
		file := LocalFiles[match[1]]
		for _, i := range pieces {
			if i >= BufferBitSize(*file) || !ByteArrayCheck(file.BufferMap.BitSequence, i) {
				fmt.Println("Invalid pieces' numbers :", i)
				return false, GetPiecesData{}
			}
		}
		return true, GetPiecesData{Key: match[1], Pieces: pieces}
	}
	return false, GetPiecesData{}
}

func (f *File) GetFile() (string, int, int, string, bool) {
	if f.Name == "" && f.Size == 0 {
		return f.Name, f.Size, f.PieceSize, f.Key, false
	}
	return f.Name, f.Size, f.PieceSize, f.Key, true
}

// Fonction de mise à jour des peers
