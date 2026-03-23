package config

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig() *AppConfig {
	cfg := NewAppConfig()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("GRPC_SERVER_PORT", 50051)

	if err := viper.ReadInConfig(); err != nil {
		logrus.Warnf("Config file not found, using defaults: %v", err)
	}

	cfg.GrpcServer.GRPCServerPort = viper.GetInt("GRPC_SERVER_PORT")
	cfg.TelegramServer.ApiId = viper.GetInt("TELEGRAM_API_ID")
	cfg.TelegramServer.ApiHash = viper.GetString("TELEGRAM_API_HASH")

	return cfg
}
