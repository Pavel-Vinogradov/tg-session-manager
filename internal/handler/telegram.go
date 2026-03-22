package handler

import (
	"context"
	"tg-session-manager/api/proto"
	"tg-session-manager/internal/config"

	"google.golang.org/grpc"
)

type TelegramHandler struct {
	proto.UnimplementedTelegramServiceServer
	TelegramService *config.TelegramService
}

func NewTelegramHandler(telegramService *config.TelegramService) *TelegramHandler {
	return &TelegramHandler{
		TelegramService: telegramService,
	}
}

func (s *TelegramHandler) CreateSession(ctx context.Context, req *proto.CreateSessionRequest) (*proto.CreateSessionResponse, error) {
	// TODO: Implement CreateSession
	return &proto.CreateSessionResponse{}, nil
}

func (s *TelegramHandler) DeleteSession(ctx context.Context, req *proto.DeleteSessionRequest) (*proto.DeleteSessionResponse, error) {
	// TODO: Implement DeleteSession
	return &proto.DeleteSessionResponse{}, nil
}

func (s *TelegramHandler) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	// TODO: Implement SendMessage
	return &proto.SendMessageResponse{}, nil
}

func (s *TelegramHandler) SubscribeMessages(req *proto.SubscribeMessagesRequest, stream grpc.ServerStreamingServer[proto.MessageUpdate]) error {
	// TODO: Implement SubscribeMessages
	return nil
}
