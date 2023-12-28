package history

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path"
	"time"
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
	slog.Info("Reading history file", "data_dir", dataDir, "history_file", historyFile)

	history := Movies{}

	if err := ensureDirExists(dataDir); err != nil {
		slog.Error("Failed to create directory", "error", err)
		return history
	}

	file, err := os.ReadFile(path.Join(dataDir, historyFile))
	if err != nil {
		slog.Error("Failed to read history file", "error", err)
		return history
	}

	if err := json.Unmarshal(file, &history); err != nil {
		slog.Error("Failed to unmarshal history file", "error", err)
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
	slog.Info("Writing history file", "data_dir", dataDir, "history_file", historyFile)

	for k, v := range data {
		history[k] = v
	}

	if err := ensureDirExists(dataDir); err != nil {
		slog.Info("Failed to create directory", "error", err)
		return
	}

	jsonString, err := json.Marshal(history)
	if err != nil {
		slog.Error("Failed to marshal history", "error", err)
		return
	}

	if err := os.WriteFile(path.Join(dataDir, historyFile), jsonString, os.ModePerm); err != nil {
		slog.Error("Failed to write history file", "error", err)
	}
}
