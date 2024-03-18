package main

import (
	"fmt"
	tools "peerproject/tools"
)

func main() {
	success, interestedData := tools.InterestedCheck("interested Uizhsja8hzUizhsja8hzUizhsja8hzsu")
	fmt.Println(success, interestedData)
	success2, interestedData2 := tools.InterestedCheck("interested izsja8hzUizhsja8hzUizhsja8hzsu")
	fmt.Println(success2, interestedData2)

	success3, haveData := tools.HaveCheck("have Uizhsja8hzUizhsja8hzUizhsja8hzsu 010010101001")
	fmt.Println(success3, haveData)
	tools.PrintBuffer(haveData.BufferMap.BitSequence)

}
