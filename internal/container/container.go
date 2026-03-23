package container

import (
	"tg-session-manager/internal/config"
	"tg-session-manager/internal/interfaces/telegram"
	"tg-session-manager/internal/server"
	"tg-session-manager/internal/service"
)

type Container struct {
	Config      *config.AppConfig
	TelegramSvc telegram.Service
	GrpcSrv     *server.GrpcServer
}

func NewContainer(cfg *config.AppConfig) *Container {
	telegramSvc := service.NewTelegramService(cfg.TelegramServer)
	grpcSrv := server.NewGrpcServer(cfg.GrpcServer)

	return &Container{
		Config:      cfg,
		TelegramSvc: telegramSvc,
		GrpcSrv:     grpcSrv,
	}
}
