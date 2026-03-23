package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"tg-session-manager/cmd/cli"
	"tg-session-manager/internal/config"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	_ = godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logrus.Info("shutdown signal received")
		cancel()
	}()

	app, err := cli.NewApp(config.LoadConfig())
	if err != nil {
		logrus.Fatalf("failed to create app: %v", err)
	}
	server := app.RegisterServiceServer()

	app.RunGrpc(server)

	<-ctx.Done()
}
