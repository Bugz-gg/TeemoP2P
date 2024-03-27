package test

import (
	"fmt"
	"peerproject/tools"
	"testing"
)

var dummyFile = tools.File{Key: "Uizhsja8hzUizhsja8hzUizhsja8hzsu", Size: 24, PieceSize: 2}
var dummyFile2 = tools.File{Key: "Uizhsja8hzUizhsja8hzU7zhsja8hzsu", Size: 36, PieceSize: 3}
var dummyFile3 = tools.File{Key: "Uizhsja8hzpolisja8hzUizhsja8hzsu", Size: 45, PieceSize: 3}

func TestList(t *testing.T) {
	fmt.Println(">>> List regex")
	tools.AddFile(tools.LocalFiles, &dummyFile)
	tools.InitBufferMap(&dummyFile)
	tools.BufferMapWrite(&dummyFile.BufferMap, 0)
	tools.BufferMapWrite(&dummyFile.BufferMap, 5)

	success, listData := tools.ListCheck("list [fe 12 1 Uizhsja8hzUizhsja8hzUizhsja8hzsu]")
	expectedListData := tools.ListData{Files: []tools.File{{Name: "fe", Size: 12, PieceSize: 1, Key: "Uizhsja8hzUizhsja8hzUizhsja8hzsu"}}}
	tools.InitBufferMap(&(expectedListData.Files[0]))
	if !success || !tools.ListDataCmp(listData, expectedListData) {
		t.Errorf("ListCheck failed. Expected: true %v, Got: %v %v", expectedListData, success, listData)
	}
}

func TestInterested(t *testing.T) {
	fmt.Println(">>> Interested regex")

	success, interestedData := tools.InterestedCheck("interested Uizhsja8hzUizhsja8hzUizhsja8hzsu")
	expectedInterestedData := tools.InterestedData{Key: "Uizhsja8hzUizhsja8hzUizhsja8hzsu"}
	if !success || !tools.InterestedCmp(interestedData, expectedInterestedData) {
		t.Errorf("InterestedCheck failed. Expected: true %v, Got: %v %v", expectedInterestedData, success, interestedData)
	}
	success2, interestedData2 := tools.InterestedCheck("interested izsja8hzUizhsja8hzUizhsja8hzsu")
	expectedInterestedData2 := tools.InterestedData{}
	if !success || !tools.InterestedCmp(interestedData2, expectedInterestedData2) {
		t.Errorf("InterestedCheck failed. Expected: true %v, Got: %v %v", expectedInterestedData2, success2, interestedData2)
	}
}

func TestHave(t *testing.T) {
	fmt.Println(">>> Have regex")

	tools.AddFile(tools.LocalFiles, &dummyFile2)
	tools.InitBufferMap(&dummyFile2)
	tools.BufferMapWrite(&dummyFile2.BufferMap, 1)
	tools.BufferMapWrite(&dummyFile2.BufferMap, 4)
	tools.BufferMapWrite(&dummyFile2.BufferMap, 6)
	tools.BufferMapWrite(&dummyFile2.BufferMap, 8)
	tools.BufferMapWrite(&dummyFile2.BufferMap, 11)

	success, haveData := tools.HaveCheck("have Uizhsja8hzUizhsja8hzU7zhsja8hzsu 010010101001")
	expectedHaveData := tools.HaveData{Key: "Uizhsja8hzUizhsja8hzU7zhsja8hzsu", BufferMap: tools.LocalFiles["Uizhsja8hzUizhsja8hzU7zhsja8hzsu"].BufferMap}
	if !success || !tools.HaveCmp(haveData, expectedHaveData) {
		t.Errorf("HaveCheck failed. Expected: true %v, Got: %v %v", expectedHaveData, success, haveData)
	}
}

func TestGetPieces(t *testing.T) {
	fmt.Println(">>> GetPieces regex")

	tools.AddFile(tools.LocalFiles, &dummyFile)
	tools.InitBufferMap(&dummyFile)
	tools.BufferMapWrite(&dummyFile.BufferMap, 0)

	tools.AddFile(tools.LocalFiles, &dummyFile3)
	tools.InitBufferMap(&dummyFile3)
	tools.BufferMapWrite(&dummyFile3.BufferMap, 0)
	tools.BufferMapWrite(&dummyFile3.BufferMap, 5)

	success, getPiecesData := tools.GetPiecesCheck("getpieces UizhsjakhzUizhsja8hzUizhsja8hzsu []")
	expectedGetPiecesData := tools.GetPiecesData{}
	if success || !tools.GetPiecesCmp(getPiecesData, expectedGetPiecesData) {
		t.Errorf("HaveCheck failed. Expected: false %v, Got: %v %v", expectedGetPiecesData, success, getPiecesData)
	}

	success2, getPiecesData2 := tools.GetPiecesCheck("getpieces Uizhsja8hzUizhsja8hzUizhsja8hzsu [0 893 88]")
	if success2 || !tools.GetPiecesCmp(getPiecesData2, expectedGetPiecesData) {
		t.Errorf("HaveCheck failed. Expected: false %v, Got: %v %v", expectedGetPiecesData, success2, getPiecesData2)
	}

	success3, getPiecesData3 := tools.GetPiecesCheck("getpieces Uizhsja8hzpolisja8hzUizhsja8hzsu [0 5]")
	expectedGetPiecesData3 := tools.GetPiecesData{Key: "Uizhsja8hzpolisja8hzUizhsja8hzsu", Pieces: []int{0, 5}}
	if !success3 || !tools.GetPiecesCmp(getPiecesData3, expectedGetPiecesData3) {
		t.Errorf("HaveCheck failed. Expected: true %v, Got: %v %v", expectedGetPiecesData3, success3, getPiecesData3)
	}
}

func TestData(t *testing.T) {
	fmt.Println(">>> Data regex")

}
