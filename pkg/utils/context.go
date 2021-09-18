package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	DefaultContext = Context{
		Endpoint:       "https://api.dns-platform.jp/dpf/v1",
		ClientLogLevel: 0,
	}
)

type Context struct {
	Endpoint       string
	Token          string
	ClientLogLevel int
}

func GetContexts() (*Context, error) {
	currentContext := viper.GetString("current-context")
	if viper.GetString("context") != "" {
		currentContext = viper.GetString("context")
	}
	contexts := viper.Get("contexts").(map[string]Context)
	c, ok := contexts[currentContext]
	if !ok {
		return nil, fmt.Errorf("failed to get config, current-context: `%s`", currentContext)
	}
	return &c, nil
}

func init() {
	viper.SetDefault("contexts", map[string]Context{"default": DefaultContext})
	viper.SetDefault("current-context", "default")
}
