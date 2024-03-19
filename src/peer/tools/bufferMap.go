package tools

import (
	"bytes"
	"fmt"
	"strconv"
)

type BufferMap struct {
	Length      int
	BitSequence []byte
}

func BufferBitSize(file File) int {
	return (file.Size-1)/file.PieceSize + 1
}

func BufferSize(file File) int {
	return (file.Size-1)/file.PieceSize/8 + 1
}

func InitBufferMap(file *File) {
	file.BufferMap = BufferMap{Length: BufferSize(*file), BitSequence: make([]byte, BufferSize(*file))}
}

func ByteArrayWrite(array []byte, index int) {
	array[index/8] |= 1 << (7 - (index % 8))
}

func ByteArrayErase(array []byte, index int) {
	array[index/8] &= ^(1 << (7 - (index % 8)))
}

func ByteArrayCheck(array []byte, index int) bool {
	return array[index/8]&(1<<(7-index%8)) > 0
}

func BufferMapWrite(bufferMap BufferMap, index int) {
	ByteArrayWrite(bufferMap.BitSequence, index)
}

func StringToBufferMap(str string) BufferMap {
	array := make([]byte, (len(str)-1)/8+1)
	for index, char := range str {
		if char == '1' {
			ByteArrayWrite(array, index)
		}
	}
	return BufferMap{Length: len(str), BitSequence: array}
}

func BufferMapToString(bufferMap BufferMap) {

}

func PrintBuffer(array []byte) {
	var buf bytes.Buffer
	for _, b := range array {
		binary := strconv.FormatInt(int64(b), 2)
		paddedBinary := fmt.Sprintf("%08s", binary)
		buf.WriteString(paddedBinary)
	}
	fmt.Println(buf.String())
}
