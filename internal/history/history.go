package history

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"time"

	log "github.com/sirupsen/logrus"
)

type Movie struct {
	Title      string    `json:"title"`
	Link       string    `json:"link"`
	CoverImage string    `json:"cover_image"`
	Year       string    `json:"year"`
	TimeStamp  time.Time `json:"timestamp"`
}

type Movies map[string]Movie

func ensureDirExists(dirPath string) error {
	if _, err := os.Stat(dirPath); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func ReadHistory(dataDir, historyFile string) Movies {
	log.Infof("Reading history file %v, %v\n", dataDir, historyFile)

	history := Movies{}

	if err := ensureDirExists(dataDir); err != nil {
		log.Error("Failed to create directory:", err)
		return history
	}

	file, err := os.ReadFile(path.Join(dataDir, historyFile))
	if err != nil {
		log.Error("Failed to read history file:", err)
		return history
	}

	if err := json.Unmarshal(file, &history); err != nil {
		log.Error("Failed to unmarshal history file:", err)
	}

	return history
}

func DiffHistory(data, history Movies) Movies {
	result := Movies{}

	for k, v := range data {
		if _, ok := history[k]; !ok {
			result[k] = v
		}
	}

	return result
}

func WriteHistory(data, history Movies, dataDir, historyFile string) {
	if len(data) == 0 {
		return
	}
	log.Infof("Writing history file %v, %v\n", dataDir, historyFile)

	for k, v := range data {
		history[k] = v
	}

	if err := ensureDirExists(dataDir); err != nil {
		log.Error("Failed to create directory:", err)
		return
	}

	jsonString, err := json.Marshal(history)
	if err != nil {
		log.Error("Failed to marshal history:", err)
		return
	}

	if err := os.WriteFile(path.Join(dataDir, historyFile), jsonString, os.ModePerm); err != nil {
		log.Error("Failed to write history file:", err)
	}
}
