package config

type (
	AppConfig struct {
		GrpcServer     *GrpcConfig
		TelegramServer *TelegramConfig
	}
)

func NewAppConfig() *AppConfig {
	return &AppConfig{
		GrpcServer:     NewGrpcConfig(),
		TelegramServer: NewTelegramConfig(),
	}
}
