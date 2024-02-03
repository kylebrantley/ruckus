package config

import (
	"fmt"
	"net/http"

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

func (c Config) Validate() []error {
	var errs []error

	err := c.Request.validateMethod()
	if err != nil {
		errs = append(errs, err)
	}

	if c.Request.URL == "" {
		errs = append(errs, fmt.Errorf("url is required"))
	}

	if c.Request.Method == "" {
		errs = append(errs, fmt.Errorf("method is required"))
	}

	if c.NumberOfRequests < 1 {
		errs = append(errs, fmt.Errorf("numberOfRequests must be greater than 0"))
	}

	if c.MaxConcurrentRequests < 1 {
		errs = append(errs, fmt.Errorf("maxConcurrentRequests must be greater than 0"))
	}

	return errs
}

func (r Request) validateMethod() error {
	allowedMethods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}

	for _, method := range allowedMethods {
		if r.Method == method {
			return nil
		}
	}

	return fmt.Errorf("invalid method: %s", r.Method)
}
