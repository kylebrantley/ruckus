package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Request               Request `mapstructure:"request"`
	NumberOfRequests      int     `mapstructure:"numberOfRequests"`
	MaxConcurrentRequests int     `mapstructure:"maxConcurrentRequests"`
}

type Request struct {
	URL     string            `mapstructure:"url"`
	Method  string            `mapstructure:"method"`
	Body    string            `mapstructure:"body"`
	Headers map[string]string `mapstructure:"headers"`
	Timeout int               `mapstructure:"timeout"`
}

func New(filePath string) (Config, error) {
	var config Config

	viper.SetConfigFile(filePath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("failed to read config file: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config file: %v", err)
	}

	return config, nil
}
