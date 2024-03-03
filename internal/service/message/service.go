package messageservice

import (
	"github.com/a1exCross/chat-server/internal/client/db"
	"github.com/a1exCross/chat-server/internal/repository"
	"github.com/a1exCross/chat-server/internal/service"
)

type serv struct {
	messageRepo repository.MessagesRepository
	txManager   db.TxManager
	logsRepo    repository.LogsRepository
}

// NewService - создает сервисный слой для сообщений
func NewService(messageRepo repository.MessagesRepository, tx db.TxManager, logsRepo repository.LogsRepository) service.MessageService {
	return &serv{
		messageRepo: messageRepo,
		txManager:   tx,
		logsRepo:    logsRepo,
	}
}
