package app

import (
	"context"
	"flag"
	"log"

	"github.com/a1exCross/chat-server/internal/config"
	"google.golang.org/grpc"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	return a, nil
}

func (a *App) Run() error {

}

func (a *App) iniDeps() error {

}

func (a *App) initConfig() error {
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load environments: %v", err)
	}

	return nil
}
