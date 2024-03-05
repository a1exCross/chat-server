package repository

import (
	"context"

	"github.com/a1exCross/chat-server/internal/model"
)

// ChatRepository - описывает методы репозитория чатов
type ChatRepository interface {
	Create(context.Context, model.ChatDTO) (int64, error)
	Delete(context.Context, int64) error
}

// MessagesRepository - описывает методы репозитория сообщений
type MessagesRepository interface {
	Create(context.Context, model.MessageDTO) error
}

// LogsRepository - описывает методы репозитория логов
type LogsRepository interface {
	Create(ctx context.Context, log model.Log) (int64, error)
	Get(ctx context.Context, id int64) (model.Log, error)
}
