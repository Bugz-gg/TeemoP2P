package tools

import (
	"fmt"
	"peerproject/tools"
	"testing"
)

var dummyFile = tools.File{Key: "Uizhsja8hzUizhsja8hzUizhsja8hzsu", Size: 24, PieceSize: 2}

func TestList(t *testing.T) {
	fmt.Println(">>> List regex")
	tools.AddFile(tools.LocalFiles, dummyFile)
	tools.InitBufferMap(&dummyFile)
	fmt.Println(dummyFile)
	tools.BufferMapWrite(dummyFile.BufferMap, 0)
	tools.BufferMapWrite(dummyFile.BufferMap, 5)
	success, listData := tools.ListCheck("list [fe 12 1 Uizhsja8hzUizhsja8hzUizhsja8hzsu]")
	fmt.Println(success, listData)
}

func TestInterested(t *testing.T) {
	fmt.Println(">>> Interested regex")

	success, interestedData := tools.InterestedCheck("interested Uizhsja8hzUizhsja8hzUizhsja8hzsu")
	fmt.Println(success, interestedData)
	success2, interestedData2 := tools.InterestedCheck("interested izsja8hzUizhsja8hzUizhsja8hzsu")
	fmt.Println(success2, interestedData2)
}

func TestHave(t *testing.T) {
	fmt.Println(">>> Have regex")
	success3, haveData := tools.HaveCheck("have Uizhsja8hzUizhsja8hzUizhsja8hzsu 010010101001")
	fmt.Println(success3, haveData)
	tools.PrintBuffer(haveData.BufferMap.BitSequence)
}

func TestGetPieces(t *testing.T) {
	fmt.Println(">>> GetPieces regex")
	tools.AddFile(tools.LocalFiles, dummyFile)
	tools.InitBufferMap(&dummyFile)
	tools.BufferMapWrite(dummyFile.BufferMap, 0)
	tools.BufferMapWrite(dummyFile.BufferMap, 5)
	success, getPiecesData := tools.GetPiecesCheck("getpieces UizhsjakhzUizhsja8hzUizhsja8hzsu []")
	fmt.Println(success, getPiecesData)
	success2, getPiecesData2 := tools.GetPiecesCheck("getpieces Uizhsja8hzUizhsja8hzUizhsja8hzsu [0 893 88]")
	fmt.Println(success2, getPiecesData2)
	success3, getPiecesData3 := tools.GetPiecesCheck("getpieces Uizhsja8hzUizhsja8hzUizhsja8hzsu [0 5]")
	fmt.Println(success3, getPiecesData3)

}
