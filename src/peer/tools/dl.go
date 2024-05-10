package tools

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/ini.v1"
)

var config map[string]map[string]string
var configMutex sync.Mutex

func FillFilesFromConfig(conn string) map[string]*File {
	fmt.Println("\033[0;36mDetecting shared files...\033[0m")
	path := GetValueFromConfig("Peer", "path")
	os.MkdirAll(filepath.Join("./", path), os.FileMode(0777))
	files, err := searchFiles(path)
	if err != nil {
		fmt.Println(err)
		return make(map[string]*File)
	}

	return fillStruct(files, conn)
}

func searchFiles(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() != "file" {
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
	ipport := strings.Split(conn, ":")
	IP, port := ipport[0], ipport[1]
	for _, filePath := range files {
		if filepath.Base(filePath) == "manifest" {
			fd, _ := os.OpenFile(filePath, os.O_RDONLY, os.FileMode(0777))
			rd := bufio.NewReader(fd)
			tempTab := []string{"Name", "Size", "PieceSize", "Key"}
			tempMap := make(map[string]string, 4)
			for _, i := range tempTab {
				txt, _ := rd.ReadString('\n')
				tempMap[i] = strings.TrimSuffix(txt, "\n")
			}
			sz, _ := strconv.ParseUint(tempMap["Size"], 10, 64)
			psz, _ := strconv.ParseUint(tempMap["PieceSize"], 10, 64)
			fil := File{
				Name:      tempMap["Name"],
				Size:      sz,
				PieceSize: psz,
				Key:       tempMap["Key"],
			}

			txt, _ := rd.ReadBytes('\n')
			buff := string(txt)
			buffermap := StringToBufferMap(buff)
			bufferMaps[fil.Key] = &buffermap
			fil.Peers = make(map[string]*Peer)

			fil.Peers["self"] = &Peer{
				IP:         IP,
				Port:       port,
				BufferMaps: bufferMaps,
			}
			fil.Peers[conn] = fil.Peers["self"]
			result[fil.Key] = &fil

			RemoteFiles[tempMap["Key"]] = &fil
			continue

		}
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return nil
		}

		fileSize := fileInfo.Size()
		pieceSizeStr := GetValueFromConfig("Peer", "default_piece_size")
		pieceSize, err := strconv.ParseUint(pieceSizeStr, 10, 64)
		if err != nil {
			return nil
		}

		pieceSize = min(pieceSize, uint64(fileSize)) // Conversion to uint64 may cause errors.

		fil := File{
			Name:      filepath.Base(filePath),
			Size:      uint64(fileSize), // Same here.
			PieceSize: pieceSize,
			Key:       GetMD5Hash(filePath),
			Complete:  true,
		}
		buffermap := InitBufferMap(fil.Size, fil.PieceSize)
		bufferMaps[fil.Key] = &buffermap
		fil.Peers = make(map[string]*Peer)
		fil.Peers["self"] = &Peer{
			IP:         IP,
			Port:       port,
			BufferMaps: bufferMaps,
		}
		fil.Peers[conn] = fil.Peers["self"]
		for u := range BufferBitSize(fil) {
			BufferMapWrite(&(*(fil.Peers["self"].BufferMaps)[fil.Key]), u)
		}
		result[fil.Key] = &fil
	}
	return result
}

func GetValueFromConfig(section string, key string) string {
	configMutex.Lock()
	if config == nil {
		config = make(map[string]map[string]string)
	}
	if config[section] == nil {
		config[section] = make(map[string]string)
	}
	if config[section][key] != "" {
		ret := config[section][key]
		configMutex.Unlock()
		return ret
	}
	file, err := ini.Load("config.ini")
	if err != nil {
		configMutex.Unlock()
		return err.Error()
	}
	sec := file.Section(section)

	valueStr := sec.Key(key).String()
	config[section][key] = valueStr
	if valueStr == "" {
		defaultValues := map[string]map[string]string{"Peer": {
			"time_dl_rare_piece":   "6000",
			"max_concurrency":      "2",
			"max_peers":            "2000",
			"max_peers_to_connect": "5",
			"progress_value":       "2",
			"update_time":          "120",
			"timeout":              "5",
			"response_timeout":     "1",
			"max_buff_size":        "8392",
			"default_piece_size":   "2048",
			"max_message_attempts": "3",
			"path":                 "share"}}
		valueStr = defaultValues[section][key]
		config[section][key] = defaultValues[section][key]
	}
	configMutex.Unlock()
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
