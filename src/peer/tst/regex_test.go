package tools

import (
	"fmt"
	"peerproject/tools"
	"testing"
)

func TestInterested(t *testing.T) {
	success, interestedData := tools.InterestedCheck("interested Uizhsja8hzUizhsja8hzUizhsja8hzsu")
	fmt.Println(success, interestedData)
	success2, interestedData2 := tools.InterestedCheck("interested izsja8hzUizhsja8hzUizhsja8hzsu")
	fmt.Println(success2, interestedData2)
}

func TestHave(t *testing.T) {

	success3, haveData := tools.HaveCheck("have Uizhsja8hzUizhsja8hzUizhsja8hzsu 010010101001")
	fmt.Println(success3, haveData)
	tools.PrintBuffer(haveData.BufferMap.BitSequence)

}
