package service

import (
	"context"

	"github.com/a1exCross/chat-server/internal/model"
)

type ChatServive interface {
	Create(context.Context, model.ChatDTO) (int64, error)
	Delete(context.Context, int64) error
}

type MessageService interface {
	SendMessage(context.Context, model.MessageDTO) error
}
