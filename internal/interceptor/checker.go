package interceptor

import (
	"context"
	"errors"
	"strings"

	"github.com/a1exCross/chat-server/internal/client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const authPrefix = "Bearer "

// AccessChecker - верификатор доступа
type AccessChecker interface {
	AccessCheck(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

type checker struct {
	authClient client.AuthService
}

// NewAccessChecker - создание верификатора доступа
func NewAccessChecker(auth client.AuthService) AccessChecker {
	return &checker{
		authClient: auth,
	}
}

// AccessCheck - проверка доступа к ручке
func (c *checker) AccessCheck(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return nil, errors.New("invalid authorization header")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	md = metadata.New(map[string]string{"Authorization": authPrefix + accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	err := c.authClient.Check(ctx, info.FullMethod)
	if err != nil {
		return nil, errors.New("access denied")
	}

	return handler(ctx, req)
}
