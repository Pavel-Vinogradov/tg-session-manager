package server

import (
	"errors"
	"fmt"
	"tg-session-manager/internal/config"

	"google.golang.org/grpc"
)

type (
	GrpcServer struct {
		Config *config.GrpcConfig
		server *grpc.Server
	}
)

func NewGrpcServer(cfg *config.GrpcConfig) *GrpcServer {
	return &GrpcServer{
		Config: cfg,
		server: grpc.NewServer(),
	}
}

func (g *GrpcServer) Health() error {
	if g.server == nil {
		return errors.New("grpc server is nil")
	}
	return nil
}

func (g *GrpcServer) Server() *grpc.Server {
	return g.server
}

func (g *GrpcServer) Address() string {
	return fmt.Sprintf(":%d", g.Config.GRPCServerPort)
}
