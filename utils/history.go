package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func ReadHistory(dataDir, historyFile string) map[string]Movie {
	log.Printf("Reading history file %v, %v\n", dataDir, historyFile)

	history := map[string]Movie{}

	if _, err := os.Stat(dataDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dataDir, os.ModePerm)
		if err != nil {
			log.Println(err)
			return history
		}
	}

	file, err := ioutil.ReadFile(path.Join(dataDir, historyFile))
	if err != nil {
		log.Println(err)
		return history
	}

	json.Unmarshal(file, &history)
	return history
}

func DiffHistory(data, history map[string]Movie) map[string]Movie {
	result := map[string]Movie{}

	for k, v := range data {
		if _, ok := history[k]; ok {
			continue
		}
		result[k] = v
	}

	return result
}

func WriteHistory(data, history map[string]Movie, dataDir, historyFile string) {

	if len(data) == 0 {
		return
	}
	log.Printf("Writing history file %v, %v\n", dataDir, historyFile)

	for k, v := range data {
		history[k] = v
	}

	if _, err := os.Stat(dataDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dataDir, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	}

	jsonString, _ := json.Marshal(history)
	ioutil.WriteFile(path.Join(dataDir, historyFile), jsonString, os.ModePerm)
}
