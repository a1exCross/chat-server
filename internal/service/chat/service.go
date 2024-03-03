package chatService

import (
	"github.com/a1exCross/chat-server/internal/client/db"
	"github.com/a1exCross/chat-server/internal/repository"
	"github.com/a1exCross/chat-server/internal/service"
)

type serv struct {
	chatRepo  repository.ChatRepository
	txManager db.TxManager
	logsRepo  repository.LogsRepository
}

func NewService(chatRepo repository.ChatRepository, tx db.TxManager, logsRepo repository.LogsRepository) service.ChatServive {
	return &serv{
		chatRepo:  chatRepo,
		txManager: tx,
		logsRepo:  logsRepo,
	}
}
