package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

const filename = "chronicle"

var (
	ErrReadInConfig = errors.New("failed to read in configuration file")
)

type Config struct{}

func New() (*Config, error) {
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%s: %s", ErrReadInConfig, err)
	}
	return &Config{}, nil
}

func (c *Config) GetMySQLHost() string {
	return viper.GetString("mysql.host")
}

func (c *Config) GetMySQLPort() uint {
	return viper.GetUint("mysql.port")
}

func (c *Config) GetMySQLServerID() uint {
	return viper.GetUint("mysql.serverID")
}

func (c *Config) GetMySQLUser() string {
	return viper.GetString("mysql.user")
}

func (c *Config) GetMySQLPassword() string {
	return viper.GetString("mysql.password")
}

func (c *Config) GetKafkaHost() string {
	return viper.GetString("kafka.host")
}

func (c *Config) GetKafkaPort() uint {
	return viper.GetUint("kafka.port")
}
