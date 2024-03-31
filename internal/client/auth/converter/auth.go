package converter

import (
	accesspb "github.com/a1exCross/chat-server/internal/client/auth/proto"
	"github.com/a1exCross/chat-server/internal/model"
)

// UserRolesToProto - Конвертирует роль пользователя в proto
func UserRolesToProto(roles ...model.UserRole) []accesspb.UserRole {
	var converted []accesspb.UserRole

	for _, v := range roles {
		converted = append(converted, accesspb.UserRole(v))
	}

	return converted
}
