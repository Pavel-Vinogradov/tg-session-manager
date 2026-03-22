package config

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig() AppConfig {
	cfg := AppConfig{}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("GRPC_SERVER_PORT", 50051)
	viper.SetDefault("SERVER_PORT", 8080)

	if err := viper.ReadInConfig(); err != nil {
		logrus.Warnf("Config file not found, using defaults: %v", err)
	}

	cfg.GrpcServer = &GrpcConfig{
		GRPCServerPort: viper.GetInt("GRPC_SERVER_PORT"),
		ServerPort:     viper.GetInt("SERVER_PORT"),
	}

	return cfg
}
