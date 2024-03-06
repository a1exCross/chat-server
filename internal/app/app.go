package app

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/a1exCross/chat-server/internal/closer"
	"github.com/a1exCross/chat-server/internal/config"
	"github.com/a1exCross/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

// App является структурой точки входа в приложение
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

// NewApp - точка входа в приложение
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run - Запуск обработчиков событиый
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {

	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, fn := range inits {
		err := fn(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load environments: %v", err)
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	chat_v1.RegisterChatV1Server(a.grpcServer, a.serviceProvider.ChatImplemetation(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	lis, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	log.Printf("Listen and serve at %s\n", a.serviceProvider.GRPCConfig().Address())

	return a.grpcServer.Serve(lis)
}
