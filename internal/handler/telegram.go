package handler

import (
	"context"
	"tg-session-manager/api/proto"
	"tg-session-manager/internal/interfaces/telegram"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Error(codes.Internal, "failed to create session")
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

	if sessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}

	err := s.telegramService.DeleteSession(sessionId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete session")
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

	if sessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}

	if text == "" {
		return nil, status.Error(codes.InvalidArgument, "message text is required")
	}

	err := s.telegramService.SendMessage(sessionId, text)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to send message")
	}

	return &proto.SendMessageResponse{}, nil
}

func (s *TelegramHandler) SubscribeMessages(req *proto.SubscribeMessagesRequest, stream grpc.ServerStreamingServer[proto.MessageUpdate]) error {
	return nil
}
