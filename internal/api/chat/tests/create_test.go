package tests

import (
	"context"
	"errors"
	"testing"

	chatapi "github.com/a1exCross/chat-server/internal/api/chat"
	"github.com/a1exCross/chat-server/internal/model"
	"github.com/a1exCross/chat-server/internal/service"
	"github.com/a1exCross/chat-server/internal/service/mocks"
	"github.com/a1exCross/chat-server/pkg/chat_v1"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	mc := minimock.NewController(t)

	type mockChat func(mc *minimock.Controller) service.ChatServive
	type mockMessage func(mc *minimock.Controller) service.MessageService

	createChatDTO := model.ChatDTO{
		Usernames: []string{"biba", "boba"},
	}

	createReq := &chat_v1.CreateRequest{
		Usernames: []string{"biba", "boba"},
	}

	id := int64(1)

	createRes := &chat_v1.CreateResponse{
		Id: id,
	}

	tests := []struct {
		name     string
		ctx      context.Context
		err      error
		req      *chat_v1.CreateRequest
		expected *chat_v1.CreateResponse
		mockMessage
		mockChat
	}{
		{
			name:     "sucessfull test",
			ctx:      ctx,
			err:      nil,
			req:      createReq,
			expected: createRes,
			mockChat: func(mc *minimock.Controller) service.ChatServive {
				mock := mocks.NewChatServiveMock(t)
				mock.CreateMock.Expect(ctx, createChatDTO).Return(id, nil)

				return mock
			},
			mockMessage: func(mc *minimock.Controller) service.MessageService {
				mock := mocks.NewMessageServiceMock(t)

				return mock
			},
		},
		{
			name:     "failed test",
			ctx:      ctx,
			err:      errors.New("failed to create chat: error"),
			req:      createReq,
			expected: nil,
			mockChat: func(mc *minimock.Controller) service.ChatServive {
				mock := mocks.NewChatServiveMock(t)
				mock.CreateMock.Expect(ctx, createChatDTO).Return(0, errors.New("error"))

				return mock
			},
			mockMessage: func(mc *minimock.Controller) service.MessageService {
				mock := mocks.NewMessageServiceMock(t)

				return mock
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			chatService := test.mockChat(mc)
			messageService := test.mockMessage(mc)

			serv := chatapi.NewImplementation(chatService, messageService)

			res, err := serv.Create(ctx, test.req)

			require.Equal(t, test.expected, res)
			require.Equal(t, test.err, err)
		})
	}

}
