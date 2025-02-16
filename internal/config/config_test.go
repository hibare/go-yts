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
	os.Setenv("GO_YTS_SCHEDULE", "test_schedule")
	os.Setenv("GO_YTS_DATA_DIR", "/test/data/dir")
	os.Setenv("GO_YTS_NOTIFIER_DISCORD_ENABLED", "false")
	os.Setenv("GO_YTS_NOTIFIER_DISCORD_WEBHOOK", "discord_webhook_url")
	os.Setenv("GO_YTS_HTTP_REQUEST_TIMEOUT", "5s")
	os.Setenv("GO_YTS_LOG_LEVEL", commonLogger.DefaultLoggerLevel)
	os.Setenv("GO_YTS_LOG_MODE", commonLogger.DefaultLoggerMode)
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

	assert.False(t, Current.NotifierConfig.Discord.Enabled)
	assert.Equal(t, "", Current.NotifierConfig.Discord.Webhook)

	assert.Equal(t, constants.DefaultRequestTimeout, Current.HTTPConfig.RequestTimeout)

	assert.Equal(t, commonLogger.DefaultLoggerLevel, Current.LoggerConfig.Level)
	assert.Equal(t, commonLogger.DefaultLoggerMode, Current.LoggerConfig.Mode)

}
