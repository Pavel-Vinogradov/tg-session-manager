package cli

import (
	"net"
	"tg-session-manager/internal/config"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type (
	App struct {
		*config.Services
	}
)

func NewApp(servicesConfig *config.Services) (*App, error) {
	app := &App{
		Services: servicesConfig,
	}
	return app, nil

}

func (app *App) RegisterServiceServer() *grpc.Server {
	server := app.GrpcServer.Server()

	reflection.Register(server)

	return server
}

func (app *App) RunGrpc(server *grpc.Server) {
	go func() {
		listener, err := net.Listen("tcp", app.GrpcServer.Address())

		if err != nil {
			logrus.Errorf("failed to listen: %v", err)
			return
		}

		logrus.Infof("grpc server listening at %v", listener.Addr())

		if err = server.Serve(listener); err != nil {
			logrus.Errorf("failed to serve: %v", err)
		}
	}()
}
