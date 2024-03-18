package main

import (
	"fmt"
	tools "peerproject/tools"
)

func main() {
	//announceRegex, _, _ := RegexInit()
	//GetAnnounceRegex := AnnounceRegex()
	//announceRegex := GetAnnounceRegex()

	tools.AnnounceCheck("announce listen 2222 seed [fe 12 1 du]")
	fmt.Println(tools.InterestedCheck("interested Uizhsja8hzUizhsja8hzUizhsja8hzsu"))
	fmt.Println(tools.InterestedCheck("interested izsja8hzUizhsja8hzUizhsja8hzsu"))

	fmt.Println(tools.HaveCheck("have Uizhsja8hzUizhsja8hzUizhsja8hzsu 010010101001"))

}
