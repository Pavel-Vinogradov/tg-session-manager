package repository

import (
	"tg-session-manager/internal/interfaces/telegram"

	tdtelegram "github.com/gotd/td/telegram"
)

type telegramRepository struct {
	client *tdtelegram.Client
}

func NewTelegramRepository(client *tdtelegram.Client) telegram.Repository {
	return &telegramRepository{
		client: client,
	}
}

func (r *telegramRepository) CreateSession() (string, string, error) {
	return "", "", nil
}

func (r *telegramRepository) DeleteSession(sessionId string) error {
	return nil
}

func (r *telegramRepository) SendMessage(sessionId, message string) error {
	return nil
}
