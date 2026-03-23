package cli

import (
	"net"
	"tg-session-manager/api/proto"
	"tg-session-manager/internal/config"
	"tg-session-manager/internal/container"
	"tg-session-manager/internal/handler"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type (
	App struct {
		container *container.Container
	}
)

func NewApp(cfg *config.AppConfig) (*App, error) {
	app := &App{
		container: container.NewContainer(cfg),
	}
	return app, nil
}

func (app *App) RegisterServiceServer() *grpc.Server {
	server := app.container.GrpcSrv.Server()

	proto.RegisterTelegramServiceServer(server, handler.NewTelegramHandler(app.container.TelegramSvc, app.container.SessionManager))

	reflection.Register(server)

	return server
}

func (app *App) RunGrpc(server *grpc.Server) {
	go func() {
		listener, err := net.Listen("tcp", app.container.GrpcSrv.Address())

		if err != nil {
			logrus.Errorf("failed to listen: %v", err)
			return
		}

		logrus.Infof("grpc handler listening at %v", listener.Addr())

		if err = server.Serve(listener); err != nil {
			logrus.Errorf("failed to serve: %v", err)
		}
	}()
}
