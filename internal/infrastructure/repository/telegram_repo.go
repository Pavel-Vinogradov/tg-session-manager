package repository

import (
	"encoding/base64"

	"tg-session-manager/internal/interfaces/telegram"

	"github.com/google/uuid"
	tdtelegram "github.com/gotd/td/telegram"
	"github.com/skip2/go-qrcode"
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
	sessionId := uuid.New().String()
	qrCode, err := r.generateQRCode(sessionId)

	return sessionId, qrCode, err

}

func (r *telegramRepository) DeleteSession(sessionId string) error {
	return nil
}

func (r *telegramRepository) SendMessage(sessionId, message string) error {
	return nil
}

func (r *telegramRepository) generateQRCode(data string) (string, error) {
	png, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png), nil
}
