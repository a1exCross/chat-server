package client

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -o ./mocks/ -s ".go"

import (
	"context"

	"github.com/a1exCross/chat-server/internal/model"
)

// AuthService - сервис авторизации и аутентификации
type AuthService interface {
	Check(ctx context.Context, role ...model.UserRole) error
}
