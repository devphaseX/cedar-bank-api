package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DbSource        string        `mapstructure:"DBSOURCE"`
	ServerAddress   string        `mapstructure:"SERVER_ADDRESS"`
	SymmetricKey    string        `mapstructure:"SYMMETRIC_KEY"`
	AccessTokenTime time.Duration `mapstructure:"ACCESS_TOKEN_TIME"`
}

func LoadConfig(path string) (config *Config, err error) {
	vp := viper.New()
	vp.AddConfigPath(path)
	vp.SetConfigName("app")
	vp.SetConfigType("env")

	vp.AutomaticEnv()
	if err = vp.ReadInConfig(); err != nil {
		return nil, err
	}

	err = vp.Unmarshal(&config)
	return
}
