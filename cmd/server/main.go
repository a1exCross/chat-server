package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	sq "github.com/Masterminds/squirrel"
	"github.com/a1exCross/chat-server/internal/config"
	pbChat "github.com/a1exCross/chat-server/pkg/chat_v1"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type service struct {
	pbChat.UnimplementedChatV1Server
	pool *pgxpool.Pool
}

const tableChats = "chats"
const tableMessages = "messages"

const (
	id        = "id"
	usernames = "usernames"
	author    = "author"
	content   = "content"
	createdAt = "created_at"
	chatID    = "chat_id"
)

func (s service) Create(ctx context.Context, req *pbChat.CreateRequest) (*pbChat.CreateResponse, error) {
	insertBuilder := sq.Insert(tableChats).
		PlaceholderFormat(sq.Dollar).
		Columns(usernames).
		Values(req.Usernames).
		Suffix(fmt.Sprintf("RETURNING %s", id))

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error at parse sql builder: %v", err)
	}

	var id int64

	err = s.pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error at query to database: %v", err)
	}

	return &pbChat.CreateResponse{
		Id: id,
	}, nil
}

func (s service) Delete(ctx context.Context, req *pbChat.DeleteRequest) (*empty.Empty, error) {
	deleteBuilder := sq.Delete(tableChats).
		Where(sq.Eq{id: req.Id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error at parse sql builder: %v", err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error at query to database: %v", err)
	}

	deleteBuilder = sq.Delete(tableMessages).
		Where(sq.Eq{id: req.Id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err = deleteBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error at parse sql builder: %v", err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error at query to database: %v", err)
	}

	return &empty.Empty{}, nil
}

func (s service) SendMessage(ctx context.Context, req *pbChat.SendMessageRequest) (*empty.Empty, error) {
	query, args, err := sq.Select(id).
		From(tableChats).
		OrderBy("id DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error at parse sql builder: %v", err)
	}

	var lastChatID int64

	err = s.pool.QueryRow(ctx, query, args...).Scan(&lastChatID)
	if err != nil {
		return nil, fmt.Errorf("error at query to database: %v", err)
	}

	insertBuilder := sq.Insert(tableMessages).
		PlaceholderFormat(sq.Dollar).
		Columns(chatID, author, content, createdAt).
		Values(lastChatID, req.Message.From, req.Message.Text, req.Message.Timestamp.AsTime()).
		Suffix(fmt.Sprintf("RETURNING %s", id))

	query, args, err = insertBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error at parse sql builder: %v", err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error at query to database: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load environments: %v", err)
	}

	grpcConf, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to create grpc config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConf.Address())
	if err != nil {
		log.Fatalf("failed to connect grpc server: %v", err)
	}

	log.Printf("Listen and serve at %s", grpcConf.Address())

	pgConf, err := config.NewPGConfig()

	log.Println(pgConf)

	if err != nil {
		log.Fatalf("failed to create pg config: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConf.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)

	pbChat.RegisterChatV1Server(
		s, service{
			pool: pool,
		})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc server: %v", err)
	}
}
