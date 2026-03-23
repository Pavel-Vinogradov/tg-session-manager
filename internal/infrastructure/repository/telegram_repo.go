package repository

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"tg-session-manager/internal/config"
	"tg-session-manager/internal/interfaces/telegram"
	"tg-session-manager/internal/session"

	"github.com/gotd/td/tg"
	"github.com/skip2/go-qrcode"
)

type telegramRepository struct {
	sm     *session.ManagerSession
	config *config.TelegramConfig
}

func NewTelegramRepository(sm *session.ManagerSession, cfg *config.TelegramConfig) telegram.Repository {
	return &telegramRepository{
		sm:     sm,
		config: cfg,
	}
}

func (r *telegramRepository) DeleteSession(sessionId string) error {
	deleted := r.sm.DeleteSession(sessionId)
	if !deleted {
		return telegram.ErrSessionNotFound
	}
	return nil
}

func (r *telegramRepository) SendMessage(sessionId, peer, message string) error {
	telegramSession, exists := r.sm.GetSession(sessionId)
	if !exists {
		return telegram.ErrSessionNotFound
	}
	if !telegramSession.IsActive() {
		return telegram.ErrSessionNotActive
	}

	api := telegramSession.Client.API()
	ctx := telegramSession.Ctx

	inputPeer, err := r.resolvePeer(ctx, api, peer)
	if err != nil {
		return fmt.Errorf("failed to resolve peer %s: %w", peer, err)
	}

	_, err = api.MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
		Peer:     inputPeer,
		Message:  message,
		RandomID: rand.Int63(),
	})

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (r *telegramRepository) generateQRCode(data string) (string, error) {
	png, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png), nil
}

func (r *telegramRepository) resolvePeer(ctx context.Context, api *tg.Client, peer string) (tg.InputPeerClass, error) {
	if len(peer) > 0 && peer[0] == '@' {
		username := peer[1:]
		peerObj, err := api.ContactsResolveUsername(ctx, &tg.ContactsResolveUsernameRequest{
			Username: username,
		})
		if err != nil {
			return nil, err
		}

		switch p := peerObj.Peer.(type) {
		case *tg.PeerUser:
			user, ok := peerObj.Users[p.UserID].(*tg.User)
			if !ok {
				return nil, fmt.Errorf("user not found")
			}
			return &tg.InputPeerUser{
				UserID:     p.UserID,
				AccessHash: user.AccessHash,
			}, nil
		case *tg.PeerChat:
			return &tg.InputPeerChat{
				ChatID: p.ChatID,
			}, nil
		case *tg.PeerChannel:
			channel, ok := peerObj.Chats[p.ChannelID].(*tg.Channel)
			if !ok {
				return nil, fmt.Errorf("channel not found")
			}
			return &tg.InputPeerChannel{
				ChannelID:  p.ChannelID,
				AccessHash: channel.AccessHash,
			}, nil
		}
	}

	return nil, fmt.Errorf("unsupported peer format: %s", peer)
}
