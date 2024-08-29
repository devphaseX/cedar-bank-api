package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DbSource          string        `mapstructure:"DBSOURCE"`
	HttpServerAddress string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GrpcServerAddress string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	SymmetricKey      string        `mapstructure:"SYMMETRIC_KEY"`
	AccessTokenTime   time.Duration `mapstructure:"ACCESS_TOKEN_TIME"`
	RefreshTokenTime  time.Duration `mapstructure:"REFRESH_TOKEN_TIME"`
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
