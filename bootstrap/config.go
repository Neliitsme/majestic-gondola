package bootstrap

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

type Config struct {
	PostgresUrl string `mapstructure:"POSTGRES_URL"`
}

func LoadConfig(l *slog.Logger) *Config {
	var config Config
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	err := viper.Unmarshal(&config)

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if l != nil {
		l.Info("Loaded config successfully.")
	}

	return &config
}
