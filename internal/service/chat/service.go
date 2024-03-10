package chatservice

import (
	"github.com/a1exCross/chat-server/internal/repository"
	"github.com/a1exCross/chat-server/internal/service"
	"github.com/a1exCross/common/pkg/client/db"
)

type serv struct {
	chatRepo  repository.ChatRepository
	txManager db.TxManager
	logsRepo  repository.LogsRepository
}

// NewService - создает сервисный слой для чатов
func NewService(chatRepo repository.ChatRepository, tx db.TxManager, logsRepo repository.LogsRepository) service.ChatServive {
	return &serv{
		chatRepo:  chatRepo,
		txManager: tx,
		logsRepo:  logsRepo,
	}
}
