package config

type (
	GrpcConfig struct {
		GRPCServerPort int `mapstructure:"grpc_port"`
		ServerPort     int `mapstructure:"server_port"`
	}
)

func NewGrpcConfig() *GrpcConfig {
	return new(GrpcConfig)
}
