package container

import (
	"tg-session-manager/internal/config"
	"tg-session-manager/internal/interfaces/telegram"
	"tg-session-manager/internal/server"
	"tg-session-manager/internal/service"
	"tg-session-manager/internal/session"
)

type Container struct {
	Config         *config.AppConfig
	TelegramSvc    telegram.Service
	SessionManager *session.ManagerSession
	GrpcSrv        *server.GrpcServer
}

func NewContainer(cfg *config.AppConfig) *Container {
	sessionManager := session.NewSessionManager()
	telegramSvc := service.NewTelegramService(cfg.TelegramServer, sessionManager)
	grpcSrv := server.NewGrpcServer(cfg.GrpcServer)

	return &Container{
		Config:         cfg,
		TelegramSvc:    telegramSvc,
		SessionManager: sessionManager,
		GrpcSrv:        grpcSrv,
	}
}
