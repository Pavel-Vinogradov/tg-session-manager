package handler

import (
	"context"
	"tg-session-manager/api/proto"
	"tg-session-manager/internal/interfaces/telegram"

	"google.golang.org/grpc"
)

type TelegramHandler struct {
	proto.UnimplementedTelegramServiceServer
	telegramService telegram.Service
}

func NewTelegramHandler(telegramService telegram.Service) *TelegramHandler {
	return &TelegramHandler{
		telegramService: telegramService,
	}
}

func (s *TelegramHandler) CreateSession(ctx context.Context, req *proto.CreateSessionRequest) (*proto.CreateSessionResponse, error) {
	sessionId, qrCode, err := s.telegramService.CreateSession()
	if err != nil {
		return nil, err
	}

	return &proto.CreateSessionResponse{
		SessionId: &sessionId,
		QrCode:    &qrCode,
	}, nil
}

func (s *TelegramHandler) DeleteSession(ctx context.Context, req *proto.DeleteSessionRequest) (*proto.DeleteSessionResponse, error) {
	var sessionId string
	if req.SessionId != nil {
		sessionId = *req.SessionId
	}

	err := s.telegramService.DeleteSession(sessionId)
	if err != nil {
		return nil, err
	}

	return &proto.DeleteSessionResponse{}, nil
}

func (s *TelegramHandler) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	var sessionId, text string
	if req.SessionId != nil {
		sessionId = *req.SessionId
	}
	if req.Text != nil {
		text = *req.Text
	}

	err := s.telegramService.SendMessage(sessionId, text)
	if err != nil {
		return nil, err
	}

	return &proto.SendMessageResponse{}, nil
}

func (s *TelegramHandler) SubscribeMessages(req *proto.SubscribeMessagesRequest, stream grpc.ServerStreamingServer[proto.MessageUpdate]) error {
	return nil
}
