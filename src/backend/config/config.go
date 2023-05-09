package config

import (
	"backend/cmd/modes/flags"
	"github.com/spf13/viper"
)

type Config struct {
	Postgres flags.PostgresFlags `mapstructure:"postgres"`
	Address  string              `mapstructure:"address"`
	Port     string              `mapstructure:"port"`
	LogPath  string              `mapstructure:"log_path"`
	LogLevel string              `mapstructure:"log_level"`
}

func (c *Config) ParseConfig(configFileName, pathToConfig string) error {
	v := viper.New()
	v.SetConfigName(configFileName)
	v.SetConfigType("json")
	v.AddConfigPath(pathToConfig)

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(c)
	if err != nil {
		return err
	}

	return nil
}
