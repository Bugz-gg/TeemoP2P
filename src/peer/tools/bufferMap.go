package tools

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
)

func BitSize(size int, pieceSize int) int {
	return int(math.Ceil(float64(size) / float64(pieceSize)))
}

// BufferBitSize returns the length of the bit sequence needed for the BufferMap of a File.
func BufferBitSize(file File) int {
	return BitSize(file.Size, file.PieceSize)
}

// BufferSize returns the length of the array containing the bit sequence for the BufferMap of a File.
func BufferSize(file File) int {
	return (file.Size-1)/file.PieceSize/8 + 1
}

// InitBufferMap helps to initialize a BufferMap of a File.
func InitBufferMap(size int, pieceSize int) BufferMap {
	return BufferMap{Length: BitSize(size, pieceSize), BitSequence: make([]byte, (size-1)/pieceSize/8+1)}
}

func LenInitBufferMap(length int) BufferMap {
	return BufferMap{Length: length, BitSequence: make([]byte, (length-1)/8+1)}
}

func BitCount(buff BufferMap) int {
	var count int = 0
	for i := range buff.Length {
		if ByteArrayCheck(buff.BitSequence, i) {
			count++
		}
	}
	return count
}

// ByteArrayWrite sets the bit at `index` position to 1.
func ByteArrayWrite(array *[]byte, index int) {
	(*array)[index/8] |= 1 << (7 - (index % 8))
}

// ByteArrayErase sets the bit at `index` position to 0.
func ByteArrayErase(array *[]byte, index int) {
	(*array)[index/8] &= ^(1 << (7 - (index % 8)))
}

// ByteArrayCheck tells if the bit at `index` position is set to 1.
func ByteArrayCheck(array []byte, index int) bool {
	return array[index/8]&(1<<(7-index%8)) > 0
}

// ArrayCheck tells you if the array is full of 1 or not.
// Usefull is you want to know if a file is entirely dl.
func ArrayCheck(array []byte) bool {
	for i := 0; i < len(array); i++ {
		if !ByteArrayCheck(array, i) {
			return false
		}
	}
	return true
}

// BufferMapWrite uses ByteArrayWrite to write a 1 at the `index` position.
func BufferMapWrite(bufferMap *BufferMap, index int) {
	ByteArrayWrite(&(bufferMap.BitSequence), index)
}

// BufferMapErase uses ByteArrayErase to write a 0 at the `index` position.
func BufferMapErase(bufferMap *BufferMap, index int) {
	ByteArrayErase(&(bufferMap.BitSequence), index)
}

// BufferMapCopy copies a BufferMap into another.
func BufferMapCopy(dst **BufferMap, src *BufferMap) {
	if *dst == nil {
		*dst = &BufferMap{Length: src.Length, BitSequence: make([]byte, (src.Length-1)/8+1)}
	}
	//for i := range dst.Length { // TODO
	for i := 0; i < (*dst).Length; i++ {
		if ByteArrayCheck(src.BitSequence, i) {
			ByteArrayWrite(&((*dst).BitSequence), i)
		} else {
			ByteArrayErase(&((*dst).BitSequence), i)
		}
	}
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
	return Data{Length: len(str) / 8, BitSequence: array}
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
