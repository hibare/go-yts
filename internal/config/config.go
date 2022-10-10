package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Schedule       string        `mapstructure:"SCHEDULE"`
	DataDir        string        `mapstructure:"DATA_DIR"`
	HistoryFile    string        `mapstructure:"HISTORY_FILE"`
	SlackWebhook   string        `mapstructure:"SLACK_WEBHOOK"`
	DiscordWebhook string        `mapstructure:"DISCORD_WEBHOOK"`
	Timeout        time.Duration `mapstructure:"TIMEOUT"`
}

func LoadConfig(path string) (config Config, err error) {

	viper.SetDefault("DATA_DIR", "/data")
	viper.SetDefault("HISTORY_FILE", "history.json")
	viper.SetDefault("SCHEDULE", "0 */4 * * *")
	viper.SetDefault("SLACK_WEBHOOK", "")
	viper.SetDefault("DISCORD_WEBHOOK", "")
	viper.SetDefault("TIMEOUT", 60)

	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	_ = viper.ReadInConfig()
	err = viper.Unmarshal(&config)
	return
}
