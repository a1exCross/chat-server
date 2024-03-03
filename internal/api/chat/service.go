package chatapi

import (
	"github.com/a1exCross/chat-server/internal/service"
	chatPb "github.com/a1exCross/chat-server/pkg/chat_v1"
)

// Implementation - структура, описывающая имплементацию gRPC сервера
type Implementation struct {
	chatPb.UnimplementedChatV1Server
	chatServ    service.ChatServive
	messageServ service.MessageService
}

// NewImplementation - создает новую имплементацию для gRPC сервера
func NewImplementation(chatServ service.ChatServive, messageServ service.MessageService) *Implementation {
	return &Implementation{
		chatServ:    chatServ,
		messageServ: messageServ,
	}
}
