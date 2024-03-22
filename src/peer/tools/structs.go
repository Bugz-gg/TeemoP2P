package tools

import (
	"strings"
)

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

type GetPiecesData struct {
	Key    string
	Pieces []int
}

type BufferMap struct {
	Length      int
	BitSequence []byte
}

func StrCmp(str1 string, str2 string) bool {
	return strings.Compare(str1, str2) == 0
}

func FileCmp(f1 File, f2 File) bool {
	return f1.Key == f2.Key && f1.Name == f2.Name && (f1.Size == f2.Size) && (f1.PieceSize == f2.PieceSize) && BufferMapCmp(f1.BufferMap, f2.BufferMap)
}

func ListDataCmp(lD1 ListData, lD2 ListData) bool {
	if len(lD1.Files) != len(lD2.Files) {
		return false
	}
	for i, f := range lD1.Files {
		if !FileCmp(f, lD2.Files[i]) {
			return false
		}
	}
	return true
}

func InterestedCmp(iD1 InterestedData, iD2 InterestedData) bool {
	return iD1.Key == iD2.Key
}

func HaveCmp(hD1 HaveData, hD2 HaveData) bool {
	return hD1.Key == hD2.Key && BufferMapCmp(hD1.BufferMap, hD2.BufferMap)
}

func GetPiecesCmp(gPD1 GetPiecesData, gPD2 GetPiecesData) bool {
	if gPD1.Key != gPD2.Key {
		return false
	}
	for i, p := range gPD1.Pieces {
		if p != gPD2.Pieces[i] {
			return false
		}
	}
	return true
}

func BufferMapCmp(bM1 BufferMap, bM2 BufferMap) bool {
	if bM1.Length != bM2.Length {
		return false
	}
	if bM1.Length == 0 {
		return true
	}
	for i := 0; i < bM1.Length; i++ {
		if bM1.BitSequence[i] != bM2.BitSequence[i] {
			return false
		}
	}
	return true
}
