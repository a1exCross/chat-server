package interceptor

import (
	"context"
	"errors"

	"github.com/a1exCross/chat-server/internal/api"
	"github.com/a1exCross/chat-server/internal/client"

	"google.golang.org/grpc"
)

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
	err := c.authClient.Check(ctx, api.RouteAccesses[info.FullMethod]...)
	if err != nil {
		return nil, errors.New("access denied")
	}

	return handler(ctx, req)
}
