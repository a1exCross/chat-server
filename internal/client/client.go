package client

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -o ./mocks/ -s ".go"

import (
	"context"
)

// AuthService - сервис авторизации и аутентификации
type AuthService interface {
	Check(ctx context.Context, endpoint string) error
}
