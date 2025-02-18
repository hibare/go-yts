package config

import (
	"log"
	"time"

	"github.com/hibare/GoCommon/v2/pkg/env"
	commonLogger "github.com/hibare/GoCommon/v2/pkg/logger"
	"github.com/hibare/go-yts/internal/constants"
)

type HTTPConfig struct {
	RequestTimeout time.Duration
}

type NotifierDiscordConfig struct {
	Enabled bool
	Webhook string
}

type NotifierConfig struct {
	Discord NotifierDiscordConfig
}

type StorageConfig struct {
	DataDir string
}

type LoggerConfig struct {
	Level string
	Mode  string
}

type Config struct {
	Schedule       string
	StorageConfig  StorageConfig
	NotifierConfig NotifierConfig
	HTTPConfig     HTTPConfig
	LoggerConfig   LoggerConfig
}

var Current *Config

func LoadConfig() {
	env.Load()
	Current = &Config{
		Schedule: env.MustString("GO_YTS_SCHEDULE", constants.DefaultSchedule),
		StorageConfig: StorageConfig{
			DataDir: env.MustString("GO_YTS_DATA_DIR", constants.DefaultDataDir),
		},
		NotifierConfig: NotifierConfig{
			Discord: NotifierDiscordConfig{
				Enabled: env.MustBool("GO_YTS_NOTIFIER_DISCORD_ENABLED", false),
				Webhook: env.MustString("GO_YTS_NOTIFIER_DISCORD_WEBHOOK", ""),
			},
		},
		HTTPConfig: HTTPConfig{
			RequestTimeout: env.MustDuration("GO_YTS_HTTP_REQUEST_TIMEOUT", constants.DefaultRequestTimeout),
		},
		LoggerConfig: LoggerConfig{
			Level: env.MustString("GO_YTS_LOG_LEVEL", commonLogger.DefaultLoggerLevel),
			Mode:  env.MustString("GO_YTS_LOG_MODE", commonLogger.DefaultLoggerMode),
		},
	}

	if !commonLogger.IsValidLogLevel(Current.LoggerConfig.Level) {
		log.Fatalf("Error invalid logger level %s", Current.LoggerConfig.Level)
	}

	if !commonLogger.IsValidLogMode(Current.LoggerConfig.Mode) {
		log.Fatalf("Error invalid logger mode %s", Current.LoggerConfig.Mode)
	}

	commonLogger.InitLogger(&Current.LoggerConfig.Level, &Current.LoggerConfig.Mode)

	if Current.NotifierConfig.Discord.Webhook == "" {
		Current.NotifierConfig.Discord.Enabled = false
	}
}
