package config

import (
	my_flags "backend/cmd/flags"
	"github.com/spf13/viper"
)

type Config struct {
	Postgres my_flags.PostgresFlags `mapstructure:"postgres"`
	Mongo    my_flags.MongoFlags    `mapstructure:"mongo"`
	Address  string                 `mapstructure:"address"`
	Port     string                 `mapstructure:"port"`
	LogPath  string                 `mapstructure:"log_path"`
	LogLevel string                 `mapstructure:"log_level"`
	Db       string                 `mapstructure:"db"`
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
