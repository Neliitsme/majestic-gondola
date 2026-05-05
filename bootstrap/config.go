package bootstrap

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresUrl             string `mapstructure:"POSTGRES_URL" validate:"required"`
	Host                    string `mapstructure:"HOST" validate:"required"`
	Port                    int    `mapstructure:"PORT" validate:"required"`
	LogLevel                string `mapstructure:"LOG_LEVEL"`
	ReviewProcessorSchedule string `mapstructure:"REVIEW_PROCESSOR_SCHEDULE"`
	ArtistProcessorSchedule string `mapstructure:"ARTIST_PROCESSOR_SCHEDULE"`
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
