package config

import (
	"fmt"
	"time"
)

type (
	AppConfig struct {
		GrpcServer *GrpcConfig
	}
	AppService struct {
		Config AppConfig
	}
	AppServiceError string

	ConnectionSettings struct {
		Timeout  time.Duration `mapstructure:"timeout"`
		Attempts int           `mapstructure:"attempts"`
	}
	ResponseSettings struct {
		Timeout time.Duration `mapstructure:"timeout"`
	}
)

func NewAppService() *ServiceChainMember {
	return &ServiceChainMember{&AppService{Config: LoadConfig()}}
}

func (a *AppService) Health() error {
	return nil
}

func (a *AppService) setupMember(services *Services) error {
	services.App = a
	return nil
}

func (a *AppService) appendMemberToServices(s *Services) {
	s.App = a
}

func (str AppServiceError) Error() string {
	return fmt.Sprintf("app service error: %s", string(str))
}
