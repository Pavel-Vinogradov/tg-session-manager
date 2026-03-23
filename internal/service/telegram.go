package service

import (
	"tg-session-manager/internal/config"
	"tg-session-manager/internal/infrastructure/repository"
	"tg-session-manager/internal/interfaces/telegram"

	tdtelegram "github.com/gotd/td/telegram"
)

type (
	TelegramService struct {
		client *tdtelegram.Client
		config *config.TelegramConfig
		repo   telegram.Repository
	}
)

func NewTelegramService(cfg *config.TelegramConfig) telegram.Service {
	client := tdtelegram.NewClient(
		cfg.ApiId,
		cfg.ApiHash,
		tdtelegram.Options{},
	)

	return &TelegramService{
		client: client,
		config: cfg,
		repo:   repository.NewTelegramRepository(client),
	}
}

func (t *TelegramService) CreateSession() (string, string, error) {
	return t.repo.CreateSession()
}

func (t *TelegramService) DeleteSession(sessionId string) error {
	return t.repo.DeleteSession(sessionId)
}

func (t *TelegramService) SendMessage(sessionId, message string) error {
	return t.repo.SendMessage(sessionId, message)
}

func (t *TelegramService) GetClient() *tdtelegram.Client {
	return t.client
}
