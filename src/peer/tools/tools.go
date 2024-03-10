// package tools
package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	//"os"
	//peer "peerproject/pair"
)

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func notZero(i int) int {
	if i != 0 {
		return 1
	}
	return 0
}

type File struct {
	name      string
	size      int
	pieceSize int
	key       string
	bufferMap BufferMap
}

type BufferMap struct {
	length      int
	bitSequence []byte
}

type announceData struct {
	port      int
	files     []File
	leechKeys []string
}

func RegexInit() (*regexp.Regexp, *regexp.Regexp, *regexp.Regexp) {
	announcePattern := `^announce\s+listen\s+(\d+)\s+seed\s+\[(.*)\]$` // Add optional leech if necessary
	announceRegex := regexp.MustCompile(announcePattern)

	lookPattern := `^look\s+\[[\w+(?:<|<=|!=|=|>|>=)\".*\"\s+]*\]$`
	lookRegex := regexp.MustCompile(lookPattern)

	getfilePattern := `^getfile\s+[a-z0-9]{32}$`
	getfileRegex := regexp.MustCompile(getfilePattern)

	//str := "announce listen 2222 seed [_]"
	//str := "look [filename=\"file_a.dat\" filesize>\"1048576\"]"
	//str := "getfile 8905e92afeb80fc7722ec89efbkdsf66"
	//match := announceRegex.FindAllStringSubmatch()
	return announceRegex, lookRegex, getfileRegex
}

func announceCheck(regex *regexp.Regexp, message string) (bool, announceData) {
	if match := regex.FindStringSubmatch(message); match != nil {
		//captured := len(match)-1
		port, _ := strconv.Atoi(match[1])
		filesData := strings.Split(match[2], " ")
		if len(filesData)%4 != 0 {
			fmt.Println("Invalid received message.")
			return false, announceData{}
		}
		announceStruct := announceData{port: port}
		nbFiles := len(filesData) / 4
		for i := 0; i < nbFiles; i++ {
			filename := filesData[i*4]
			size, err := strconv.Atoi(filesData[i*4+1])
			pieceSize, err2 := strconv.Atoi(filesData[i*4+2])

			if err != nil || err2 != nil {
				fmt.Println("Invalid conversion to int (size or piece size).", err, err2)
				return false, announceData{}
			}
			key := filesData[i*4+3]
			//namePattern := `[a-zA-Z_][a-zA-Z0-9_]*`
			if len(key) != 32 {
				errors.New("Key error.")
			}
			announceStruct.files = append(announceStruct.files, File{name: filename, size: size, pieceSize: pieceSize, key: key, bufferMap: BufferMap{length: (size-1)/pieceSize/8 + 1, bitSequence: make([]byte, (size-1)/pieceSize/8+1)}})
		}
		fmt.Println(announceStruct.files)

		// Prendre en compte l'absence de fichier et leech

	}
	return false, announceData{}
}

func main() {
	announceRegex, _, _ := RegexInit()

	announceCheck(announceRegex, "announce listen 2222 seed [fe 12 1 du]")
}

// Fonction de mise à jour des peers
