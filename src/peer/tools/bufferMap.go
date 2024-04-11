package tools

import (
	"bytes"
	"fmt"
	"strconv"
)

// BufferBitSize returns the length of the bit sequence needed for the BufferMap of a File.
func BufferBitSize(file File) int {
	return file.Size / file.PieceSize
}

// BufferSize returns the length of the array containing the bit sequence for the BufferMap of a File.
func BufferSize(file File) int {
	return (file.Size-1)/file.PieceSize/8 + 1
}

// InitBufferMap helps to initialize a BufferMap of a File.
func InitBufferMap(file *File) {
	file.BufferMap = BufferMap{Length: BufferBitSize(*file), BitSequence: make([]byte, BufferSize(*file))}
}

// ByteArrayWrite sets the bit at `index` position to 1.
func ByteArrayWrite(array *[]byte, index int) {
	(*array)[index/8] |= 1 << (7 - (index % 8))
}

// ByteArrayErase sets the bit at `index` position to 0.
func ByteArrayErase(array []byte, index int) {
	array[index/8] &= ^(1 << (7 - (index % 8)))
}

// ByteArrayCheck tells if the bit at `index` position is set to 1.
func ByteArrayCheck(array []byte, index int) bool {
	return array[index/8]&(1<<(7-index%8)) > 0
}

// BufferMapWrite uses ByteArrayWrite to write a 1 at the `index` position.
func BufferMapWrite(bufferMap *BufferMap, index int) {
	ByteArrayWrite(&(bufferMap.BitSequence), index)
}

// StringToBufferMap transforms a string of `0` and `1` into a BufferMap.
func StringToBufferMap(str string) BufferMap {
	array := make([]byte, (len(str)-1)/8+1)
	for index, char := range str {
		if char == '1' {
			ByteArrayWrite(&array, index)
		}
	}
	return BufferMap{Length: len(str), BitSequence: array}
}

// StringToData transforms a string of `0` and `1` into a Data.
func StringToData(str string) Data {
	array := make([]byte, (len(str)-1)/8+1)
	for index, char := range str {
		if char == '1' {
			ByteArrayWrite(&array, index)
		}
	}
	return Data{Length: len(str), BitSequence: array}
}

// BufferMapToString transforms a BufferMap into a string of `0` and `1`.
func BufferMapToString(bufferMap BufferMap) string {
	var buf bytes.Buffer
	for _, b := range bufferMap.BitSequence {
		binary := strconv.FormatInt(int64(b), 2)
		paddedBinary := fmt.Sprintf("%08s", binary)
		buf.WriteString(paddedBinary)
	}
	return buf.String()[:bufferMap.Length]

}

// PrintBuffer helps debugging by printing the content of an array of bytes.
func PrintBuffer(array []byte) {
	var buf bytes.Buffer
	for _, b := range array {
		binary := strconv.FormatInt(int64(b), 2)
		paddedBinary := fmt.Sprintf("%08s", binary)
		buf.WriteString(paddedBinary)
	}
	fmt.Println(buf.String())
}

// WriteFile writes the received data for a file.
func WriteFile(file *File, index int, str string) {

}
