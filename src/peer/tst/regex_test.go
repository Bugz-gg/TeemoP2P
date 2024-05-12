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
	tmpFiles := map[string]*tools.File{}
	tools.LocalFiles = &tmpFiles
	tools.AddFile(tools.LocalFiles, &dummyFile)
	//tools.InitBufferMap(&dummyFile)
	//tools.BufferMapWrite(&dummyFile.BufferMap, 0)
	//tools.BufferMapWrite(&dummyFile.BufferMap, 5)

	success, listData := tools.ListCheck("list [f.e 12 1 Uizhsja8hzUizhsjp8hzUizhsja8hzsu]\r\n")
	expectedListData := tools.ListData{Files: []tools.File{{Name: "f.e", Size: 12, PieceSize: 1, Key: "Uizhsja8hzUizhsjp8hzUizhsja8hzsu"}}}
	if !success || !tools.ListDataCmp(listData, expectedListData) {
		t.Errorf("ListCheck failed. Expected: true %v, Got: %v %v", expectedListData, success, listData)
	}

	success2, listData2 := tools.ListCheck("list [etoausk 22 2 df833476b1fbb8aa113c14d5a9421180]\n")
	expectedListData2 := tools.ListData{Files: []tools.File{{Name: "etoausk", Size: 22, PieceSize: 2, Key: "df833476b1fbb8aa113c14d5a9421180"}}}
	if !success || !tools.ListDataCmp(listData2, expectedListData2) {
		t.Errorf("ListCheck failed. Expected: true %v, Got: %v %v", expectedListData2, success2, listData2)
	}

	success3, listData3 := tools.ListCheck("list [etoausk 22 2 df8334l6b1fbb8aa113c14d5a9421180]\r\n")
	expectedListData3 := tools.ListData{Files: []tools.File{{Name: "etoausk(1)", Size: 22, PieceSize: 2, Key: "df8334l6b1fbb8aa113c14d5a9421180"}}}
	if !success || !tools.ListDataCmp(listData3, expectedListData3) {
		t.Errorf("ListCheck failed. Expected: true %v, Got: %v %v", expectedListData3, success3, listData3)
	}

	success4, listData4 := tools.ListCheck("list [f.e 12 1 Uizhsja8hzUizhsjm8hzUizhsja8hzsu]\n")
	expectedListData4 := tools.ListData{Files: []tools.File{{Name: "f.e(1)", Size: 12, PieceSize: 1, Key: "Uizhsja8hzUizhsjm8hzUizhsja8hzsu"}}}
	if !success || !tools.ListDataCmp(listData4, expectedListData4) {
		t.Errorf("ListCheck failed. Expected: true %v, Got: %v %v", expectedListData4, success4, listData4)
	}

	success5, listData5 := tools.ListCheck("list [etoausk 22 2 df8334l6b1fbn8aa113c14d5a9421180]\r\n")
	expectedListData5 := tools.ListData{Files: []tools.File{{Name: "etoausk(2)", Size: 22, PieceSize: 2, Key: "df8334l6b1fbn8aa113c14d5a9421180"}}}
	if !success || !tools.ListDataCmp(listData5, expectedListData5) {
		t.Errorf("ListCheck failed. Expected: true %v, Got: %v %v", expectedListData5, success5, listData5)
	}

	success6, listData6 := tools.ListCheck("list [etoausk 22 2 df833476b1fbb8aa113c14d5a9421180]\n")
	expectedListData6 := tools.ListData{Files: []tools.File{{Name: "etoausk", Size: 22, PieceSize: 2, Key: "df833476b1fbb8aa113c14d5a9421180"}}}
	if !success || !tools.ListDataCmp(listData6, expectedListData6) {
		t.Errorf("ListCheck failed. Expected: true %v, Got: %v %v", expectedListData6, success6, listData6)
	}
}

func TestInterested(t *testing.T) {
	fmt.Println(">>> Interested regex")

	success, interestedData := tools.InterestedCheck("interested Uizhsja8hzUizhsja8hzUizhsja8hzsu\r\n")
	expectedInterestedData := tools.InterestedData{Key: "Uizhsja8hzUizhsja8hzUizhsja8hzsu"}
	if !success || !tools.InterestedCmp(interestedData, expectedInterestedData) {
		t.Errorf("InterestedCheck failed. Expected: true %v, Got: %v %v", expectedInterestedData, success, interestedData)
	}
	success2, interestedData2 := tools.InterestedCheck("interested izsja8hzUizhsja8hzUizhsja8hzsu\n")
	expectedInterestedData2 := tools.InterestedData{}
	if !success || !tools.InterestedCmp(interestedData2, expectedInterestedData2) {
		t.Errorf("InterestedCheck failed. Expected: true %v, Got: %v %v", expectedInterestedData2, success2, interestedData2)
	}
}

func TestHave(t *testing.T) {
	fmt.Println(">>> Have regex")
	tmpFiles := map[string]*tools.File{}
	tools.LocalFiles = &tmpFiles
	tmpFilesRemote := map[string]*tools.File{}
	tools.RemoteFiles = tmpFilesRemote
	tmpMap := make(map[string]*tools.BufferMap)
	tmpMap[dummyFile2.Key] = &tools.BufferMap{Length: 12, BitSequence: make([]byte, 2)}
	tools.ByteArrayWrite(&tmpMap[dummyFile2.Key].BitSequence, 1)
	tools.ByteArrayWrite(&tmpMap[dummyFile2.Key].BitSequence, 4)
	tools.ByteArrayWrite(&tmpMap[dummyFile2.Key].BitSequence, 6)
	tools.ByteArrayWrite(&tmpMap[dummyFile2.Key].BitSequence, 8)
	tools.ByteArrayWrite(&tmpMap[dummyFile2.Key].BitSequence, 11)
	//dummyPeer := tools.Peer{IP: "10.0.0.1", Port: 34, BufferMaps: &tmpMap}
	tools.AddFile(&tools.RemoteFiles, &dummyFile2)
	//tools.InitBufferMap(&dummyFile2)

	success, haveData := tools.HaveCheck("have Uizhsja8hzUizhsja8hzU7zhsja8hzsu 010010101001\r\n")
	expectedHaveData := tools.HaveData{Key: "Uizhsja8hzUizhsja8hzU7zhsja8hzsu", BufferMap: tools.BufferMap{Length: 12, BitSequence: tmpMap[dummyFile2.Key].BitSequence}}
	if !success || !tools.HaveCmp(haveData, expectedHaveData) {
		t.Errorf("HaveCheck failed. Expected: true %v, Got: %v %v", expectedHaveData, success, haveData)
	}
}

func TestGetPieces(t *testing.T) {
	fmt.Println(">>> GetPieces regex")

	tmpFiles := map[string]*tools.File{}
	tools.LocalFiles = &tmpFiles
	tools.AddFile(tools.LocalFiles, &dummyFile)
	//tools.InitBufferMap(&dummyFile)
	tmpMap := make(map[string]*tools.BufferMap)
	tmpMap[dummyFile.Key] = &tools.BufferMap{Length: 12, BitSequence: make([]byte, 2)}
	tools.ByteArrayWrite(&tmpMap[dummyFile.Key].BitSequence, 0)
	tools.AddFile(tools.LocalFiles, &dummyFile)
	dummyBuffermap := tools.InitBufferMap(dummyFile.Size, dummyFile.PieceSize)
	dummyBufferMaps := make(map[string]*tools.BufferMap)
	dummyBufferMaps[dummyFile.Key] = &dummyBuffermap
	dummyFile.Peers = make(map[string]*tools.Peer)
	dummyFile.Peers["self"] = &tools.Peer{BufferMaps: dummyBufferMaps}
	tools.ByteArrayWrite(&dummyFile.Peers["self"].BufferMaps["Uizhsja8hzUizhsja8hzUizhsja8hzsu"].BitSequence, 0)
	// tools.BufferMapWrite(&dummyFile.BufferMap, 0)

	tools.AddFile(tools.LocalFiles, &dummyFile3)
	dummyBuffermap2 := tools.InitBufferMap(dummyFile3.Size, dummyFile3.PieceSize)
	dummyBufferMaps2 := make(map[string]*tools.BufferMap)
	dummyBufferMaps2[dummyFile3.Key] = &dummyBuffermap2
	dummyFile3.Peers = make(map[string]*tools.Peer)
	dummyFile3.Peers["self"] = &tools.Peer{BufferMaps: dummyBufferMaps2}

	//tools.InitBufferMap(&dummyFile3)

	tmpMap[dummyFile3.Key] = &tools.BufferMap{Length: 12, BitSequence: make([]byte, 2)}
	tools.ByteArrayWrite(&dummyFile3.Peers["self"].BufferMaps["Uizhsja8hzpolisja8hzUizhsja8hzsu"].BitSequence, 0)
	tools.ByteArrayWrite(&dummyFile3.Peers["self"].BufferMaps["Uizhsja8hzpolisja8hzUizhsja8hzsu"].BitSequence, 5)

	success, getPiecesData := tools.GetPiecesCheck("getpieces UizhsjakhzUizhsja8hzUizhsja8hzsu []\r\n")
	expectedGetPiecesData := tools.GetPiecesData{}
	if success || !tools.GetPiecesCmp(getPiecesData, expectedGetPiecesData) {
		t.Errorf("GetPiecesCheck failed. Expected: false %v, Got: %v %v", expectedGetPiecesData, success, getPiecesData)
	}
	success2, getPiecesData2 := tools.GetPiecesCheck("getpieces Uizhsja8hzUizhsja8hzUizhsja8hzsu [0 893 88]\n")
	expectedGetPiecesData2 := tools.GetPiecesData{Key: "Uizhsja8hzUizhsja8hzUizhsja8hzsu", Pieces: []int{0}}
	if !success2 || !tools.GetPiecesCmp(getPiecesData2, expectedGetPiecesData2) {
		t.Errorf("GetPiecesCheck failed. Expected: true %v, Got: %v %v", expectedGetPiecesData2, success2, getPiecesData2)
	}

	success3, getPiecesData3 := tools.GetPiecesCheck("getpieces Uizhsja8hzpolisja8hzUizhsja8hzsu [0 5]\r\n")
	expectedGetPiecesData3 := tools.GetPiecesData{Key: "Uizhsja8hzpolisja8hzUizhsja8hzsu", Pieces: []int{0, 5}}
	if !success3 || !tools.GetPiecesCmp(getPiecesData3, expectedGetPiecesData3) {
		t.Errorf("GetPiecesCheck failed. Expected: true %v, Got: %v %v", expectedGetPiecesData3, success3, getPiecesData3)
	}

	tools.ByteArrayErase(&dummyFile.Peers["self"].BufferMaps["Uizhsja8hzUizhsja8hzUizhsja8hzsu"].BitSequence, 0)
	tools.ByteArrayErase(&dummyFile3.Peers["self"].BufferMaps["Uizhsja8hzpolisja8hzUizhsja8hzsu"].BitSequence, 0)
	tools.ByteArrayErase(&dummyFile3.Peers["self"].BufferMaps["Uizhsja8hzpolisja8hzUizhsja8hzsu"].BitSequence, 5)
}

func TestData(t *testing.T) {
	fmt.Println(">>> Data regex")

	//tools.AddFile(tools.LocalFiles, &dummyFile)
	tmpMap := make(map[string]*tools.BufferMap)
	tmpMap[dummyFile.Key] = &tools.BufferMap{Length: 12, BitSequence: make([]byte, 2)}
	tools.ByteArrayWrite(&tmpMap[dummyFile.Key].BitSequence, 0)
	//tools.InitBufferMap(&dummyFile)
	//tools.BufferMapWrite(&dummyFile.BufferMap, 0)

	tmpMap[dummyFile3.Key] = &tools.BufferMap{Length: 15, BitSequence: make([]byte, 2)}
	tools.ByteArrayWrite(&tmpMap[dummyFile3.Key].BitSequence, 0)
	tmpMap[dummyFile3.Key] = &tools.BufferMap{Length: 12, BitSequence: make([]byte, 2)}
	tools.ByteArrayWrite(&tmpMap[dummyFile3.Key].BitSequence, 0)
	tools.ByteArrayWrite(&tmpMap[dummyFile3.Key].BitSequence, 5)
	tools.AddFile(&tools.RemoteFiles, &dummyFile3)
	//tools.InitBufferMap(&dummyFile3)
	//tools.BufferMapWrite(&dummyFile3.BufferMap, 0)
	//tools.BufferMapWrite(&dummyFile3.BufferMap, 5)

	// No piece given
	success, dataData := tools.DataCheck("data UizhsjakhzUizhsja8hzUizhsja8hzsu []\n")
	expectedDataData := tools.DataData{}
	if success || !tools.DataCmp(dataData, expectedDataData) {
		t.Errorf("DataCheck failed. Expected: false %v, Got: %v %v", expectedDataData, success, dataData)
	}

	// Wrong piece size
	success2, dataData2 := tools.DataCheck("data Uizhsja8hzUizhsja8hzUizhsja8hzsu [0:0 893:0 88:1]\r\n")
	if success2 || !tools.DataCmp(dataData2, expectedDataData) {
		t.Errorf("DataCheck failed. Expected: false %v, Got: %v %v", expectedDataData, success2, dataData2)
	}

	// Out of range piece index
	success3, dataData3 := tools.DataCheck("data Uizhsja8hzUizhsja8hzUizhsja8hzsu [0:0101101101101010 893:1010011001100110 88:1011110010101010]\n")
	if success2 || !tools.DataCmp(dataData3, expectedDataData) {
		t.Errorf("DataCheck failed. Expected: false %v, Got: %v %v", expectedDataData, success3, dataData3)
	}

	// Correct
	success4, dataData4 := tools.DataCheck("data Uizhsja8hzpolisja8hzUizhsja8hzsu [0:lpm 5:dsj]\r\n")
	expectedDataData4 := tools.DataData{Key: "Uizhsja8hzpolisja8hzUizhsja8hzsu", Pieces: []tools.Piece{{Index: 0, Data: tools.Data{String: []byte("lpm")}}, {Index: 5, Data: tools.Data{String: []byte("dsj")}}}}
	if !success4 || !tools.DataCmp(dataData4, expectedDataData4) {
		t.Errorf("DataCheck failed. Expected: true %v, Got: %v %v", expectedDataData4, success4, dataData4)
	}

	// Correct
	success5, dataData5 := tools.DataCheck("data Uizhsja8hzpolisja8hzUizhsja8hzsu [0:5k1]\n")
	expectedDataData5 := tools.DataData{Key: "Uizhsja8hzpolisja8hzUizhsja8hzsu", Pieces: []tools.Piece{{Index: 0, Data: tools.Data{String: []byte("5k1")}}}}
	if !success5 || !tools.DataCmp(dataData5, expectedDataData5) {
		t.Errorf("DataCheck failed. Expected: true %v, Got: %v %v", expectedDataData5, success5, dataData5)
	}
}

func TestPeers(t *testing.T) {
	fmt.Println(">>> Peers regex")
	tmpFiles := map[string]*tools.File{}
	tools.LocalFiles = &tmpFiles
	tools.AddFile(tools.LocalFiles, &dummyFile)

	tools.AllPeers["10.0.0.10:32"] = &tools.Peer{IP: "10.0.0.10", Port: "32"}
	tools.AllPeers["249.111.109.19:100"] = &tools.Peer{IP: "249.111.109.19", Port: "100"}

	tools.RemoteFiles["UizhsjakhzUizhsja8hzUizhsja8hzsu"] = &tools.File{Key: "UizhsjakhzUizhsja8hzUizhsja8hzsu"}
	// No peer given
	success := tools.PeersCheck("peers UizhsjakhzUizhsja8hzUizhsja8hzsu []\n", "")
	expectedPeersMap := make(map[string]*tools.Peer)
	if success || !tools.MapPeersCmp(tools.RemoteFiles["UizhsjakhzUizhsja8hzUizhsja8hzsu"].Peers, expectedPeersMap) {
		t.Errorf("PeersCheck failed. Expected: false (%v), Got: %v (%v)", expectedPeersMap, success, (*tools.LocalFiles)["UizhsjakhzUizhsja8hzUizhsja8hzsu"].Peers)
	}

	tools.RemoteFiles["Uizhsja8hzUizhsja8hzUizhsja8hzsu"] = &tools.File{Key: "Uizhsja8hzUizhsja8hzUizhsja8hzsu"}
	// Wrong peer format
	success2 := tools.PeersCheck("peers Uizhsja8hzUizhsja8hzUizhsja8hzsu [0:0 893:0 88:1]\r\n", "")
	if success2 || !tools.MapPeersCmp(tools.RemoteFiles["Uizhsja8hzUizhsja8hzUizhsja8hzsu"].Peers, expectedPeersMap) {
		t.Errorf("PeersCheck failed. Expected: false (%v), Got: %v (%v)", expectedPeersMap, success2, (*tools.LocalFiles)["Uizhsja8hzUizhsja8hzUizhsja8hzsu"].Peers)
	}

	tools.RemoteFiles["Uizhsja8hzpolisja8hzUizhsja8hzsu"] = &tools.File{Key: "Uizhsja8hzUizhsja8hzUizhsja8hzsu"}
	tools.RemoteFiles["Uizhsja8hzpolisja8hzUizhsja8hzsu"].Peers = map[string]*tools.Peer{}
	// Correct
	success4 := tools.PeersCheck("peers Uizhsja8hzpolisja8hzUizhsja8hzsu [10.0.0.10:32 249.111.109.19:100]\n", "")
	expectedPeersMap4 := make(map[string]map[string]*tools.Peer)
	expectedPeersMap4["Uizhsja8hzpolisja8hzUizhsja8hzsu"] = map[string]*tools.Peer{"10.0.0.10:32": {IP: "10.0.0.10", Port: "32"}, "249.111.109.19:100": {IP: "249.111.109.19", Port: "100"}}
	//expectedPeersData4 := tools.PeersData{Key: "Uizhsja8hzpolisja8hzUizhsja8hzsu", Peers: []tools.Peer{{IP: "10.0.0.10", Port: 32}, {IP: "249.111.109.19", Port: 100}}}
	if !success4 || !tools.MapPeersCmp(tools.RemoteFiles["Uizhsja8hzUizhsja8hzUizhsja8hzsu"].Peers, expectedPeersMap4["Uizhsja8hzpolisja8hzUizhsja8hzsu"]) {
		t.Errorf("PeersCheck failed. Expected: true (%v), Got: %v (%v)", expectedPeersMap4, success4, (*tools.LocalFiles)["Uizhsja8hzUizhsja8hzUizhsja8hzsu"].Peers)
	}
}
