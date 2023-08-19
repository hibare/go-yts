package history

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistoryReadWrite(t *testing.T) {
	tempDir := t.TempDir() // Creates a temporary directory for testing

	dataDir := filepath.Join(tempDir, "data")
	historyFile := "history.json"

	mockData := Movies{
		"movie1": {Title: "Movie 1"},
		"movie2": {Title: "Movie 2"},
	}

	defer os.RemoveAll(tempDir) // Clean up the temporary directory

	WriteHistory(mockData, make(Movies), dataDir, historyFile)
	history := ReadHistory(dataDir, historyFile)

	assert.Equal(t, mockData, history)
}

func TestDiffHistory(t *testing.T) {
	data := Movies{
		"movie1": {Title: "Movie 1"},
		"movie2": {Title: "Movie 2"},
		"movie3": {Title: "Movie 3"},
	}

	history := Movies{
		"movie1": {Title: "Movie 1"},
		"movie2": {Title: "Movie 2"},
	}

	expectedDiff := Movies{
		"movie3": {Title: "Movie 3"},
	}

	diff := DiffHistory(data, history)
	assert.Equal(t, expectedDiff, diff)
}

func TestEmptyDataDirectory(t *testing.T) {
	tempDir := t.TempDir()
	dataDir := filepath.Join(tempDir, "data")
	historyFile := "history.json"

	// Read history from an empty data directory
	history := ReadHistory(dataDir, historyFile)
	assert.Empty(t, history)
}

func TestUnreadableHistoryFile(t *testing.T) {
	tempDir := t.TempDir()
	dataDir := filepath.Join(tempDir, "data")
	historyFile := "history.json"

	// Create an unreadable history file
	unreadableFilePath := filepath.Join(dataDir, historyFile)
	file, _ := os.Create(unreadableFilePath)
	file.Chmod(0000)
	file.Close()

	// Attempt to read history from the unreadable file
	history := ReadHistory(dataDir, historyFile)
	assert.Empty(t, history)
}

func TestWriteHistoryWithNoData(t *testing.T) {
	tempDir := t.TempDir()
	dataDir := filepath.Join(tempDir, "data")
	historyFile := "history.json"

	// Write history with no data
	WriteHistory(make(Movies), make(Movies), dataDir, historyFile)

	// Verify that no file was written
	_, err := os.Stat(filepath.Join(dataDir, historyFile))
	assert.True(t, os.IsNotExist(err))
}

func TestMarshalFail(t *testing.T) {
	tempDir := t.TempDir() // Creates a temporary directory for testing

	dataDir := filepath.Join(tempDir, "data")
	historyFile := "history_marshal_fail.json"

	defer os.RemoveAll(tempDir) // Clean up the temporary directory

	// Create a mock data that should not cause Marshal failure
	mockData := Movies{
		"movie1":        {Title: "Movie 1"},
		"invalid_movie": {}, // Should not cause Marshal failure
	}

	WriteHistory(mockData, make(Movies), dataDir, historyFile)

	history := ReadHistory(dataDir, historyFile)
	assert.NotEmpty(t, history)
}

func TestUnmarshalFail(t *testing.T) {
	tempDir := t.TempDir() // Creates a temporary directory for testing

	dataDir := filepath.Join(tempDir, "data")
	historyFile := "history_unmarshal_fail.json"

	defer os.RemoveAll(tempDir) // Clean up the temporary directory

	// Create a history file with invalid JSON content
	invalidJSON := []byte(`{"invalid_movie":}`)
	invalidFilePath := filepath.Join(dataDir, historyFile)
	file, _ := os.Create(invalidFilePath)
	file.Write(invalidJSON)
	file.Close()

	history := ReadHistory(dataDir, historyFile)
	assert.Empty(t, history)
}
