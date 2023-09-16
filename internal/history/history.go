package history

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog/log"
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
	log.Info().Msgf("Reading history file %s, %s", dataDir, historyFile)

	history := Movies{}

	if err := ensureDirExists(dataDir); err != nil {
		log.Error().Err(err).Msg("Failed to create directory")
		return history
	}

	file, err := os.ReadFile(path.Join(dataDir, historyFile))
	if err != nil {
		log.Error().Err(err).Msg("Failed to read history file")
		return history
	}

	if err := json.Unmarshal(file, &history); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal history file")
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
	log.Info().Msgf("Writing history file %s, %s\n", dataDir, historyFile)

	for k, v := range data {
		history[k] = v
	}

	if err := ensureDirExists(dataDir); err != nil {
		log.Error().Err(err).Msg("Failed to create directory")
		return
	}

	jsonString, err := json.Marshal(history)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal history")
		return
	}

	if err := os.WriteFile(path.Join(dataDir, historyFile), jsonString, os.ModePerm); err != nil {
		log.Error().Err(err).Msg("Failed to write history file")
	}
}
