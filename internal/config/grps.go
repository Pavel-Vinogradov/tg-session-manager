package config

import (
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
