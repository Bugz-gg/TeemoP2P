package tools

import (
	"regexp"
)

type File struct {
	Name      string
	Size      int
	PieceSize int
	Key       string
	BufferMap BufferMap
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

func (f *File) GetFile() (string, int, int, string, bool) {
	if f.Name == "" && f.Size == 0 {
		return f.Name, f.Size, f.PieceSize, f.Key, false
	}
	return f.Name, f.Size, f.PieceSize, f.Key, true
}

// Fonction de mise à jour des peers
