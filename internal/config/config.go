package config

import (
	"time"

	"github.com/hibare/GoCommon/v2/pkg/env"
	"github.com/hibare/go-yts/internal/constants"
)

type HTTPConfig struct {
	RequestTimeout time.Duration
}

type NotifierSlackConfig struct {
	Enabled bool
	Webhook string
}

type NotifierDiscordConfig struct {
	Enabled bool
	Webhook string
}

type NotifierConfig struct {
	Slack   NotifierSlackConfig
	Discord NotifierDiscordConfig
}

type StorageConfig struct {
	DataDir     string
	HistoryFile string
}

type Config struct {
	Schedule       string
	StorageConfig  StorageConfig
	NotifierConfig NotifierConfig
	HTTPConfig     HTTPConfig
}

var Current *Config

func LoadConfig() {
	env.Load()
	Current = &Config{
		Schedule: env.MustString("SCHEDULE", constants.DefaultSchedule),
		StorageConfig: StorageConfig{
			DataDir:     env.MustString("DATA_DIR", constants.DefaultDataDir),
			HistoryFile: env.MustString("HISTORY_FILE", constants.DefaultHistoryFilename),
		},
		NotifierConfig: NotifierConfig{
			Slack: NotifierSlackConfig{
				Enabled: env.MustBool("NOTIFIER_SLACK_ENABLED", false),
				Webhook: env.MustString("NOTIFIER_SLACK_WEBHOOK", ""),
			},
			Discord: NotifierDiscordConfig{
				Enabled: env.MustBool("NOTIFIER_DISCORD_ENABLED", false),
				Webhook: env.MustString("NOTIFIER_DISCORD_WEBHOOK", ""),
			},
		},
		HTTPConfig: HTTPConfig{
			RequestTimeout: env.MustDuration("HTTP_REQUEST_TIMEOUT", constants.DefaultRequestTimeout),
		},
	}

	if Current.NotifierConfig.Discord.Webhook == "" {
		Current.NotifierConfig.Discord.Enabled = false
	}

	if Current.NotifierConfig.Slack.Webhook == "" {
		Current.NotifierConfig.Slack.Enabled = false
	}
}
