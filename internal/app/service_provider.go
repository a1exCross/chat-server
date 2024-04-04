package app

import (
	"context"
	"log"

	chatAPI "github.com/a1exCross/chat-server/internal/api/chat"
	"github.com/a1exCross/chat-server/internal/client"
	"github.com/a1exCross/chat-server/internal/client/auth"
	"github.com/a1exCross/chat-server/internal/config"
	"github.com/a1exCross/chat-server/internal/interceptor"
	"github.com/a1exCross/chat-server/internal/repository"
	chatRepo "github.com/a1exCross/chat-server/internal/repository/chat"
	logsRepo "github.com/a1exCross/chat-server/internal/repository/log"
	messagesRepo "github.com/a1exCross/chat-server/internal/repository/message"
	"github.com/a1exCross/chat-server/internal/service"
	chatService "github.com/a1exCross/chat-server/internal/service/chat"
	messageService "github.com/a1exCross/chat-server/internal/service/message"

	"github.com/a1exCross/common/pkg/client/db"
	"github.com/a1exCross/common/pkg/client/db/pg"
	"github.com/a1exCross/common/pkg/client/db/transaction"
	"github.com/a1exCross/common/pkg/closer"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.HTTPConfig
	authConfig    config.AuthConfig

	dbClient  db.Client
	txManager db.TxManager

	chatsRepo    repository.ChatRepository
	messagesRepo repository.MessagesRepository
	logsRepo     repository.LogsRepository

	chatService    service.ChatService
	messageService service.MessageService

	chatImpl *chatAPI.Implementation

	authClient client.AuthService

	accessChecker interceptor.AccessChecker
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get gRPC config: %v", err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) AuthConfig() config.AuthConfig {
	if s.authConfig == nil {
		cfg, err := config.NewAuthConfig()
		if err != nil {
			log.Fatalf("failed to get auth config: %v", err)
		}

		s.authConfig = cfg
	}

	return s.authConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %v", err)
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %v", err)
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %v", err)
		}

		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatsRepo == nil {
		s.chatsRepo = chatRepo.NewRepository(s.DBClient(ctx))
	}

	return s.chatsRepo
}

func (s *serviceProvider) MessagesRepository(ctx context.Context) repository.MessagesRepository {
	if s.messagesRepo == nil {
		s.messagesRepo = messagesRepo.NewRepository(s.DBClient(ctx))
	}

	return s.messagesRepo
}

func (s *serviceProvider) LogsRepository(ctx context.Context) repository.LogsRepository {
	if s.logsRepo == nil {
		s.logsRepo = logsRepo.NewRepository(s.DBClient(ctx))
	}

	return s.logsRepo
}

func (s *serviceProvider) TxManager(_ context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.dbClient.DB())
	}

	return s.txManager
}

func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(
			s.ChatRepository(ctx),
			s.TxManager(ctx),
			s.LogsRepository(ctx),
		)
	}

	return s.chatService
}

func (s *serviceProvider) MessageService(ctx context.Context) service.MessageService {
	if s.messageService == nil {
		s.messageService = messageService.NewService(
			s.MessagesRepository(ctx),
			s.TxManager(ctx),
			s.LogsRepository(ctx),
		)
	}

	return s.messageService
}

func (s *serviceProvider) ChatImplementation(ctx context.Context) *chatAPI.Implementation {
	if s.chatImpl == nil {
		s.chatImpl = chatAPI.NewImplementation(s.ChatService(ctx), s.MessageService(ctx))
	}

	return s.chatImpl
}

func (s *serviceProvider) AuthClient(_ context.Context) client.AuthService {
	if s.authClient == nil {
		authClient, err := auth.NewAuthClient(s.AuthConfig())
		if err != nil {
			log.Fatalf("failed to connect auth service: %v", err)
		}

		s.authClient = authClient
	}

	return s.authClient
}

func (s *serviceProvider) AccessChecker(ctx context.Context) interceptor.AccessChecker {
	if s.accessChecker == nil {
		s.accessChecker = interceptor.NewAccessChecker(s.AuthClient(ctx))
	}

	return s.accessChecker
}
