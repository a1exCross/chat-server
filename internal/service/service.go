package service

import (
	"context"

	"github.com/a1exCross/chat-server/internal/model"
)

// ChatServive - интерфейс, описывающий сервисный слой чатов
type ChatServive interface {
	Create(context.Context, model.ChatDTO) (int64, error)
	Delete(context.Context, int64) error
}

// MessageService - интерфейс, описывающий сервисный слой сообщений
type MessageService interface {
	SendMessage(context.Context, model.MessageDTO) error
}
