package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"tg-session-manager/internal/config"
	"tg-session-manager/internal/infrastructure/repository"
	"tg-session-manager/internal/interfaces/telegram"
	"tg-session-manager/internal/session"

	"github.com/google/uuid"
	tdsession "github.com/gotd/td/session"
	tdtelegram "github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
)

type (
	TelegramService struct {
		config *config.TelegramConfig
		repo   telegram.Repository
		sm     *session.ManagerSession
	}
)

func NewTelegramService(cfg *config.TelegramConfig, sm *session.ManagerSession) telegram.Service {
	if sm == nil {
		sm = session.NewSessionManager()
	}

	return &TelegramService{
		config: cfg,
		repo:   repository.NewTelegramRepository(sm, cfg),
		sm:     sm,
	}
}

func (t *TelegramService) CreateSession() (string, string, error) {
	sessionId := uuid.New().String()

	if err := os.MkdirAll(t.config.SessionDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create session directory: %w", err)
	}

	sessionPath := filepath.Join(t.config.SessionDir, sessionId)
	storage := &tdsession.FileStorage{Path: sessionPath}

	dispatcher := tg.NewUpdateDispatcher()
	dispatcher.OnNewMessage(func(ctx context.Context, e tg.Entities, u *tg.UpdateNewMessage) error {
		t.handleIncomingMessage(sessionId, u)
		return nil
	})

	client := tdtelegram.NewClient(
		t.config.ApiId,
		t.config.ApiHash,
		tdtelegram.Options{
			SessionStorage: storage,
			UpdateHandler:  dispatcher,
		},
	)

	telegramSession := t.sm.CreateSession(sessionId, client)
	go t.startClient(telegramSession)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var token string
	select {
	case token = <-telegramSession.QRChan:
	case <-ctx.Done():
		return sessionId, "", fmt.Errorf("timeout waiting for QR token")
	}

	qrURL := "tg://login?token=" + token
	png, err := t.generateQRCode(qrURL)
	if err != nil {
		return sessionId, "", err
	}
	qrCode := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)

	return sessionId, qrCode, nil
}

func (t *TelegramService) DeleteSession(sessionId string) error {
	return t.repo.DeleteSession(sessionId)
}

func (t *TelegramService) SendMessage(sessionId, peer, message string) error {
	return t.repo.SendMessage(sessionId, peer, message)
}

func (t *TelegramService) GetSessionManager() *session.ManagerSession {
	return t.sm
}

func (t *TelegramService) startClient(telegramSession *session.TelegramSession) {
	err := telegramSession.Client.Run(telegramSession.Ctx, func(ctx context.Context) error {
		return t.handleAuth(ctx, telegramSession)
	})
	if err != nil {
		logrus.Errorf("Client error for session %s: %v", telegramSession.ID, err)
	}
}

func (t *TelegramService) handleAuth(ctx context.Context, telegramSession *session.TelegramSession) error {
	api := telegramSession.Client.API()

	for {
		resp, err := api.AuthExportLoginToken(ctx, &tg.AuthExportLoginTokenRequest{
			APIID:     t.config.ApiId,
			APIHash:   t.config.ApiHash,
			ExceptIDs: []int64{},
		})
		if err != nil {
			return err
		}

		switch v := resp.(type) {
		case *tg.AuthLoginToken:
			token := base64.RawURLEncoding.EncodeToString(v.Token)
			select {
			case telegramSession.QRChan <- token:
				logrus.Infof("QR token sent for session %s", telegramSession.ID)
			case <-ctx.Done():
				return ctx.Err()
			}

		case *tg.AuthLoginTokenMigrateTo:
			logrus.Infof("Migrate to DC %d for session %s", v.DCID, telegramSession.ID)
			return fmt.Errorf("migration not implemented")

		case *tg.AuthLoginTokenSuccess:
			logrus.Infof("Auth successful for session %s", telegramSession.ID)
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(25 * time.Second):
		}
	}
}

func (t *TelegramService) handleIncomingMessage(sessionID string, u *tg.UpdateNewMessage) {
	msg, ok := u.Message.(*tg.Message)
	if !ok || msg.Out || msg.Message == "" {
		return
	}

	sess, exists := t.sm.GetSession(sessionID)
	if !exists || !sess.IsActive() {
		return
	}

	from := "Unknown"
	if msg.FromID != nil {
		switch v := msg.FromID.(type) {
		case *tg.PeerUser:
			from = fmt.Sprintf("%d", v.UserID)
		case *tg.PeerChat:
			from = fmt.Sprintf("chat_%d", v.ChatID)
		case *tg.PeerChannel:
			from = fmt.Sprintf("channel_%d", v.ChannelID)
		}
	}

	update := &session.MessageUpdate{
		MessageID: int64(msg.ID),
		From:      from,
		Text:      msg.Message,
		Timestamp: int64(msg.Date),
	}

	sess.SendUpdate(update)
}

func (t *TelegramService) generateQRCode(data string) ([]byte, error) {
	png, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}
	return png, nil
}
