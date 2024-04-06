package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// File contains the basic data about a file.
type File struct {
	Name      string
	Size      int
	PieceSize int
	Key       string
	BufferMap BufferMap
}

// Piece contains the data about a piece (the actual data).
type Piece struct {
	Index int
	Data  Data
}

// Data contains the actual data of a piece.
// The `Length` attribute is the length of the bit sequence, not of the `BitSequence` array, which is padded to the byte.
// Not to be confused with the BufferMap struct which contains the bits telling whether a peer has pieces.
type Data struct {
	Length      int
	BitSequence []byte
}

// BufferMap represents the buffer map for the file.
type BufferMap struct {
	// Implementation details here
}

func readConfigFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "path=") {
			parts := strings.Split(line, "=")
			if len(parts) != 2 {
				return "", fmt.Errorf("invalid line: %s", line)
			}
			return parts[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("path not found in config file")
}

func searchFiles(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, filePath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func fillStruct(files []string) ([]File, error) {
	var result []File
	for _, filePath := range files {
		// Process file data and fill the File struct
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}

		fileSize := fileInfo.Size()

		file := File{
			Name: filepath.Base(filePath),
			Size: int(fileSize),
			// You can fill other fields as needed
		}

		result = append(result, file)
	}
	return result, nil
}

func main() {
	// Read path from config file
	path, err := readConfigFile("config.ini")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	// Search for files in the specified path
	files, err := searchFiles(path)
	if err != nil {
		fmt.Println("Error searching files:", err)
		return
	}

	// Fill the struct with file data
	fileStructs, err := fillStruct(files)
	if err != nil {
		fmt.Println("Error filling struct:", err)
		return
	}

	// Do something with the filled struct
	fmt.Println(fileStructs)
}
