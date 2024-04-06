package tools

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/ini.v1"
)

func FillStructFromConfig() []File {
	path, err := readConfigFile()
	if err != nil {
		return nil
	}

	files, err := searchFiles(path)
	if err != nil {
		return nil
	}

	return fillStruct(files)
}

func readConfigFile() (string, error) {
	file, err := ini.Load("config.ini")
	if err != nil {
		return "", err
	}

	section := file.Section("Peer")
	path := section.Key("path").String()

	return path, nil
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

func fillStruct(files []string) []File {
	var result []File
	for _, filePath := range files {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return nil
		}

		fileSize := fileInfo.Size()
		file, err := ini.Load("config.ini")
		if err != nil {
			return nil
		}
		section := file.Section("Peer")

		pieceSizeStr := section.Key(filepath.Base(filePath)).String()
		pieceSize, err := strconv.Atoi(pieceSizeStr)
		if err != nil {
			pieceSizeStr = section.Key("length_piece_default").String()
			pieceSize, err = strconv.Atoi(pieceSizeStr)
			if err != nil {
				return nil
			}
		}

		fil := File{
			Name:      filepath.Base(filePath),
			Size:      int(fileSize),
			PieceSize: pieceSize,
			Key:       GetMD5Hash(filePath),
		}
		result = append(result, fil)
	}
	return result
}

func GetMD5Hash(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}
