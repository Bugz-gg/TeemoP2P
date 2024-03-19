package tools

import (
	"fmt"
	"peerproject/tools"
	"testing"
)

func TestList(t *testing.T) {
	fmt.Println(">>> List regex")
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
