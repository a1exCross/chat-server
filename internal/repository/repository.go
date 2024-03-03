package repository

import "context"

type ChatRepository interface {
}

type MessagesRepository interface {
	Create(ctx context.Context, message string) error
}

type LogsRepository interface {
}
