package bootstrap

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresUrl string `mapstructure:"POSTGRES_URL" validate:"required"`
	Address     string `mapstructure:"ADDRESS"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
}

func LoadConfig() *Config {
	var config Config
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := validator.New().Struct(&config); err != nil {
		panic(fmt.Errorf("config validation error: %w", err))
	}

	return &config
}
