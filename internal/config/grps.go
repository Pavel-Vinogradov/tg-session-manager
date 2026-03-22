package config

import (
	"errors"
	"fmt"

	"google.golang.org/grpc"
)

type (
	GrpcServer struct {
		Config *GrpcConfig
		server *grpc.Server
	}
	GrpcConfig struct {
		GRPCServerPort int `mapstructure:"grpc_port"`
		ServerPort     int `mapstructure:"server_port"`
	}
)

func NewGrpcServer() *ServiceChainMember {
	return &ServiceChainMember{&GrpcServer{Config: &GrpcConfig{}}}
}

func (g *GrpcServer) Health() error {
	if g.Server() == nil {
		return errors.New("grpc server is nil")
	}

	return nil
}

func (g *GrpcServer) setupMember(services *Services) error {
	if g.Config == nil {
		return errors.New("grpc config is nil")
	}

	if services.App != nil && services.App.Config.GrpcServer != nil {
		g.Config.GRPCServerPort = services.App.Config.GrpcServer.GRPCServerPort
		g.Config.ServerPort = services.App.Config.GrpcServer.ServerPort
	}

	g.server = grpc.NewServer()

	return nil
}

func (g *GrpcServer) appendMemberToServices(services *Services) {
	services.GrpcServer = g

}

func (g *GrpcServer) Server() *grpc.Server {
	return g.server
}

func (g *GrpcServer) Address() string {
	return fmt.Sprintf(":%d", g.Config.GRPCServerPort)
}
