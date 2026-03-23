package handler

import (
	"context"
	"tg-session-manager/api/proto"
	"tg-session-manager/internal/interfaces/telegram"
	"tg-session-manager/internal/session"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TelegramHandler struct {
	proto.UnimplementedTelegramServiceServer
	telegramService telegram.Service
	sessionManager  *session.ManagerSession
}

func NewTelegramHandler(telegramService telegram.Service, sessionManager *session.ManagerSession) *TelegramHandler {
	return &TelegramHandler{
		telegramService: telegramService,
		sessionManager:  sessionManager,
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
	var sessionId, peer, text string
	if req.SessionId != nil {
		sessionId = *req.SessionId
	}
	if req.Peer != nil {
		peer = *req.Peer
	}
	if req.Text != nil {
		text = *req.Text
	}

	if sessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}

	if peer == "" {
		return nil, status.Error(codes.InvalidArgument, "peer is required")
	}

	if text == "" {
		return nil, status.Error(codes.InvalidArgument, "message text is required")
	}

	err := s.telegramService.SendMessage(sessionId, peer, text)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to send message")
	}

	return &proto.SendMessageResponse{}, nil
}

func (s *TelegramHandler) SubscribeMessages(req *proto.SubscribeMessagesRequest, stream grpc.ServerStreamingServer[proto.MessageUpdate]) error {
	var sessionId string
	if req.SessionId != nil {
		sessionId = *req.SessionId
	}

	if sessionId == "" {
		return status.Error(codes.InvalidArgument, "session_id is required")
	}

	telegramSession, exists := s.sessionManager.GetSession(sessionId)
	if !exists {
		return status.Error(codes.NotFound, "session not found")
	}

	if !telegramSession.IsActive() {
		return status.Error(codes.FailedPrecondition, "session is not active")
	}

	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update, ok := <-telegramSession.UpdatesCh:
			if !ok {
				return nil
			}

			err := stream.Send(&proto.MessageUpdate{
				MessageId: &update.MessageID,
				From:      &update.From,
				Text:      &update.Text,
				Timestamp: &update.Timestamp,
			})
			if err != nil {
				return err
			}
		}
	}
}
