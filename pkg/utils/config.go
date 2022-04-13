package utils

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	DefaultContext = Context{
		Endpoint:       "https://api.dns-platform.jp/dpf/v1",
		ClientLogLevel: 0,
	}
)

type Config struct {
	CurrentContext string             `yaml:"currentContext"`
	Contexts       map[string]Context `yaml:"contexts"`
}

type Context struct {
	Endpoint       string `yaml:"endpoint"`
	Token          string `yaml:"token"`
	ClientLogLevel int    `yaml:"loglevel"`
}

func GetConfig() (*Config, error) {
	c := &Config{
		CurrentContext: "default",
		Contexts:       map[string]Context{"default": DefaultContext},
	}
	if err := viper.Unmarshal(c); err != nil {
		return c, err
	}
	return c, nil
}

func (c *Config) WriteConfig() error {
	filename := viper.GetViper().ConfigFileUsed()
	if filename == "" {
		home, err := homedir.Dir()
		cobra.CheckErr(err)
		filename = path.Join(home, ".dpfctl.yaml")
	}
	bs, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	f, err := os.OpenFile(filename, flags, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open configfile %s: %w", filename, err)
	}
	defer f.Close()
	if _, err := f.Write(bs); err != nil {
		return fmt.Errorf("failed to write configfile %s: %w", filename, err)
	}
	return nil
}

func GetContexts() (*Context, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}
	currentContext := cfg.CurrentContext
	if viper.GetString("context") != "" {
		currentContext = viper.GetString("context")
	}

	c, ok := cfg.Contexts[currentContext]
	if !ok {
		return nil, fmt.Errorf("failed to get config, current-context: `%s`", currentContext)
	}
	return &c, nil
}
