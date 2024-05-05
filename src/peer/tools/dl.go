package tools

import (
	"crypto/md5"
	"encoding/hex"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var config map[string]map[string]string
var configMutex sync.Mutex

func FillFilesFromConfig(conn string) map[string]*File {
	path := GetValueFromConfig("Peer", "path")
	if path == "" {
		path = "share"
	}
	os.MkdirAll(filepath.Join("./", path), os.FileMode(0777))
	files, err := searchFiles(path)
	if err != nil {
		return nil
	}

	return fillStruct(files, conn)
}

func searchFiles(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() != "file" && info.Name() != "manifest" {
			files = append(files, filePath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func fillStruct(files []string, conn string) map[string]*File {
	result := make(map[string]*File)
	bufferMaps := make(map[string]*BufferMap)
	for _, filePath := range files {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return nil
		}

		fileSize := fileInfo.Size()
		pieceSizeStr := GetValueFromConfig("Peer", filepath.Base(filePath))
		pieceSize, err := strconv.ParseUint(pieceSizeStr, 10, 64)
		if err != nil {
			pieceSizeStr = GetValueFromConfig("Peer", "default_piece_size")
			pieceSize, err = strconv.ParseUint(pieceSizeStr, 10, 64)
			if err != nil {
				return nil
			}
		}
		pieceSize = min(pieceSize, uint64(fileSize)) // Conversion to uint64 may cause errors.

		fil := File{
			Name:      filepath.Base(filePath),
			Size:      uint64(fileSize), // Same here.
			PieceSize: pieceSize,
			Key:       GetMD5Hash(filePath),
		}
		buffermap := InitBufferMap(fil.Size, fil.PieceSize)
		bufferMaps[fil.Key] = &buffermap
		//InitBufferMap(&fil)
		if fil.Peers == nil {
			fil.Peers = make(map[string]*Peer)
		}
		fil.Peers["self"] = &Peer{
			IP:         "",
			Port:       "",
			BufferMaps: bufferMaps,
		}
		fil.Peers[conn] = fil.Peers["self"]
		// for u := range BufferSize(fil) {
		for u := range BufferBitSize(fil) {
			BufferMapWrite(&(*(fil.Peers["self"].BufferMaps)[fil.Key]), u)
		}
		// fmt.Println(fil.BufferMap)
		result[fil.Key] = &fil
	}
	return result
}

func GetValueFromConfig(section string, key string) string {
	if config == nil {
		configMutex.Lock()
		config = make(map[string]map[string]string)
		configMutex.Unlock()
	}
	if config[section] == nil {
		configMutex.Lock()
		config[section] = make(map[string]string)
		configMutex.Unlock()
	}
	if config[section][key] != "" {
		return config[section][key]
	}
	file, err := ini.Load("config.ini")
	if err != nil {
		return err.Error()
	}
	sec := file.Section(section)

	valueStr := sec.Key(key).String()
	configMutex.Lock()
	config[section][key] = valueStr
	configMutex.Unlock()
	if valueStr == "" {
		configMutex.Lock()
		defaultValues := map[string]map[string]string{"Peer": {"time_dl_rare_piece": "6000",
			"max_concurrency":      "2",
			"max_peers":            "2000",
			"max_peers_to_connect": "5",
			"progress_value":       "2",
			"update_time":          "120",
			"timeout":              "5",
			"max_buff_size":        "8392",
			"default_piece_size":   "2048",
			"max_message_attempts": "3",
			"path":                 "share"}}
		valueStr = defaultValues[section][key]
		config[section][key] = defaultValues[section][key]
		configMutex.Unlock()
	}
	return valueStr
}

func GetMD5Hash(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}

	hash := md5.New()
	if _, err := hash.Write(data); err != nil {
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}
