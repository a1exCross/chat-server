package converter

import (
	"github.com/a1exCross/chat-server/internal/model"
	chatPb "github.com/a1exCross/chat-server/pkg/chat_v1"
)

// ProtoToMessage - преобразует protobuf в messageDTO
func ProtoToMessage(message *chatPb.Message) model.MessageDTO {
	return model.MessageDTO{
		Author:    message.From,
		Content:   message.Text,
		CreatedAt: message.Timestamp.AsTime(),
	}
}
