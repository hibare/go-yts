package config

import (
	"os"
	"testing"
	"time"

	commonLogger "github.com/hibare/GoCommon/v2/pkg/logger"
	"github.com/hibare/go-yts/internal/constants"
	"github.com/stretchr/testify/assert"
)

// Mock environment variables for testing
func mockEnvVariables() {
	os.Setenv("SCHEDULE", "test_schedule")
	os.Setenv("DATA_DIR", "/test/data/dir")
	os.Setenv("HISTORY_FILE", "test_history.txt")
	os.Setenv("NOTIFIER_SLACK_ENABLED", "true")
	os.Setenv("NOTIFIER_SLACK_WEBHOOK", "slack_webhook_url")
	os.Setenv("NOTIFIER_DISCORD_ENABLED", "false")
	os.Setenv("NOTIFIER_DISCORD_WEBHOOK", "discord_webhook_url")
	os.Setenv("HTTP_REQUEST_TIMEOUT", "5s")
	os.Setenv("LOG_LEVEL", commonLogger.DefaultLoggerLevel)
	os.Setenv("LOG_MODE", commonLogger.DefaultLoggerMode)
}

func TestLoadConfig(t *testing.T) {
	mockEnvVariables()
	defer func() {
		// Clean up after the test
		os.Clearenv()
	}()

	LoadConfig()

	assert.Equal(t, "test_schedule", Current.Schedule)
	assert.Equal(t, "/test/data/dir", Current.StorageConfig.DataDir)
	assert.Equal(t, "test_history.txt", Current.StorageConfig.HistoryFile)

	assert.True(t, Current.NotifierConfig.Slack.Enabled)
	assert.Equal(t, "slack_webhook_url", Current.NotifierConfig.Slack.Webhook)

	assert.False(t, Current.NotifierConfig.Discord.Enabled)
	assert.Equal(t, "discord_webhook_url", Current.NotifierConfig.Discord.Webhook)

	assert.Equal(t, 5*time.Second, Current.HTTPConfig.RequestTimeout)

	assert.Equal(t, commonLogger.DefaultLoggerLevel, Current.LoggerConfig.Level)
	assert.Equal(t, commonLogger.DefaultLoggerMode, Current.LoggerConfig.Mode)
}

func TestLoadConfigWithDefaults(t *testing.T) {
	LoadConfig()

	assert.Equal(t, constants.DefaultSchedule, Current.Schedule)
	assert.Equal(t, constants.DefaultDataDir, Current.StorageConfig.DataDir)
	assert.Equal(t, constants.DefaultHistoryFilename, Current.StorageConfig.HistoryFile)

	assert.False(t, Current.NotifierConfig.Slack.Enabled)
	assert.Equal(t, "", Current.NotifierConfig.Slack.Webhook)

	assert.False(t, Current.NotifierConfig.Discord.Enabled)
	assert.Equal(t, "", Current.NotifierConfig.Discord.Webhook)

	assert.Equal(t, constants.DefaultRequestTimeout, Current.HTTPConfig.RequestTimeout)

	assert.Equal(t, commonLogger.DefaultLoggerLevel, Current.LoggerConfig.Level)
	assert.Equal(t, commonLogger.DefaultLoggerMode, Current.LoggerConfig.Mode)

}
