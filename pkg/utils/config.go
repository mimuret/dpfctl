package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Endpoint       string
	Token          string
	ClientLogLevel int
}

func NewConfig() *Config {
	return &Config{
		Endpoint:       "https://api.dns-platform.jp/dpf/v1",
		ClientLogLevel: 0,
	}
}

func GetConfig() (*Config, error) {
	c := NewConfig()
	if err := viper.Unmarshal(c); err != nil {
		fmt.Println("config file Unmarshal error")
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}
	return c, nil
}
