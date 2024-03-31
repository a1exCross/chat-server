package auth

import (
	"context"
	"strings"

	"github.com/a1exCross/chat-server/internal/client"
	"github.com/a1exCross/chat-server/internal/client/auth/converter"
	accesspb "github.com/a1exCross/chat-server/internal/client/auth/proto"
	"github.com/a1exCross/chat-server/internal/config"
	"github.com/a1exCross/chat-server/internal/model"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type auth struct {
	client     accesspb.AccessV1Client
	authConfig config.AuthConfig
}

const authPrefix = "Bearer "

// NewAuthClient - создает новый инстанс подключения к сервису auth
func NewAuthClient(authConfig config.AuthConfig) (client.AuthService, error) {
	conn, err := grpc.Dial(authConfig.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	cl := accesspb.NewAccessV1Client(conn)
	return &auth{
		client:     cl,
		authConfig: authConfig,
	}, nil
}

func (a *auth) Check(ctx context.Context, roles ...model.UserRole) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return errors.New("invalid authorization header")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	md = metadata.New(map[string]string{"Authorization": authPrefix + accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := a.client.Check(ctx, &accesspb.CheckRequest{
		Role: converter.UserRolesToProto(roles...),
	})

	return err
}
