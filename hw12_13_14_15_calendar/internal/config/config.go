package config

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

type Config interface {
	Default(string, interface{})
	Set(string, interface{})
	GetInt(string) int
	GetString(string) string
}

type config struct {
	v    *viper.Viper
	path string
}

func New(path string) (Config, error) {
	vc := viper.New()

	vc.AutomaticEnv()
	vc.SetEnvKeyReplacer(strings.NewReplacer(",", "_"))
	vc.SetEnvPrefix("config")

	if path == "" {
		return nil, errors.New("Config file is not specified")
	}

	vc.SetConfigFile(path)
	if err := vc.ReadInConfig(); err != nil {
		return nil, err
	}

	return &config{
		v:    vc,
		path: path,
	}, nil
}

func (c *config) Default(key string, value interface{}) {
	c.v.SetDefault(key, value)
}

func (c *config) Set(key string, value interface{}) {
	c.v.Set(key, value)
}

func (c *config) GetInt(key string) int {
	return c.v.GetInt(key)
}

func (c *config) GetString(key string) string {
	return c.v.GetString(key)
}
