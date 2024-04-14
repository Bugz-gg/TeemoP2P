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

func FillFilesFromConfig() map[string]*File {
	path := GetValueFromConfig("Peer", "path")
	if path == "" {
		return nil
	}

	files, err := searchFiles(path)
	if err != nil {
		return nil
	}

	return fillStruct(files)
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

func fillStruct(files []string) map[string]*File {
	result := make(map[string]*File)
	bufferMaps := make(map[string]*BufferMap)
	for _, filePath := range files {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return nil
		}

		fileSize := fileInfo.Size()
		pieceSizeStr := GetValueFromConfig("Peer", filepath.Base(filePath))
		pieceSize, err := strconv.Atoi(pieceSizeStr)
		if err != nil {
			pieceSizeStr = GetValueFromConfig("Peer", "length_piece_default")
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
		buffermap := InitBufferMap(fil.Size, fil.PieceSize)
		bufferMaps[fil.Key] = &buffermap
		//InitBufferMap(&fil)
		fil.Peers["self"] = &Peer{
			IP:         "",
			Port:       0,
			BufferMaps: bufferMaps,
		}
		for u := range BufferSize(fil) {
			BufferMapWrite(&(*(fil.Peers["self"].BufferMaps)[fil.Key]), u)
		}
		// fmt.Println(fil.BufferMap)
		result[fil.Key] = &fil
	}
	return result
}

func GetValueFromConfig(section string, key string) string {
	file, err := ini.Load("config.ini")
	if err != nil {
		return err.Error()
	}
	sec := file.Section(section)

	pieceSizeStr := sec.Key(key)
	return pieceSizeStr.String()

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
