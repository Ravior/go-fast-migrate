package util

import (
	"github.com/golobby/config/v2"
	"github.com/golobby/config/v2/feeder"
	"os"
)

var BasePath, _ = os.Getwd()

var ConfigHelper = &configHelper{}

type configHelper struct {
	config *config.Config
}

func (c *configHelper) init() {
	if c.config == nil {
		_config, err := config.New(&feeder.YamlDirectory{Path: BasePath + "/config"})
		SysHelper.CheckErr(err)

		c.config = _config
	}
}

func (c *configHelper) GetString(key string, _default ...string) string {
	c.init()
	value, err := c.config.GetString(key)
	if err != nil {
		defaultValue := ""
		if len(_default) > 0 {
			defaultValue = _default[0]
		}
		return defaultValue
	}
	return value
}

func (c *configHelper) GetInt(key string, _default ...int) int {
	c.init()
	value, err := c.config.GetInt(key)
	if err != nil {
		defaultValue := 0
		if len(_default) > 0 {
			defaultValue = _default[0]
		}
		return defaultValue
	}
	return value
}
