package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net"

	descChat "github.com/a1exCross/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	grpcPort = 50052
)

type server struct {
	descChat.UnimplementedChatV1Server
}

// Chat - структура, описывающая чатик
type Chat struct {
	ID    int64
	Users []string
}

var chat Chat

func (s server) Create(_ context.Context, req *descChat.CreateRequest) (*descChat.CreateResponse, error) {
	if len(req.Usernames) != 0 {
		id, err := rand.Int(rand.Reader, big.NewInt(123))
		if err != nil {
			return nil, fmt.Errorf("failed to generate id: %v", err)
		}

		chat.ID = id.Int64()
		chat.Users = req.Usernames

		log.Printf("chat was created with id %d", chat.ID)

		return &descChat.CreateResponse{
			Id: chat.ID,
		}, nil
	}

	return nil, fmt.Errorf("uesrnames does not exist")
}

func (s server) Delete(_ context.Context, req *descChat.DeleteRequest) (*emptypb.Empty, error) {
	if req.Id == chat.ID {
		log.Printf("chat was deleted with id %d", chat.ID)
		chat = Chat{}

		return &emptypb.Empty{}, nil
	}

	return nil, fmt.Errorf("chat not found")
}

func (s server) SendMessage(_ context.Context, req *descChat.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("message exist at %v from %s: %s", req.Message.Timestamp.AsTime(), req.Message.From, req.Message.Text)
	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", grpcPort, err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	descChat.RegisterChatV1Server(s, server{})

	log.Printf("server listening at: %d", grpcPort)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}
