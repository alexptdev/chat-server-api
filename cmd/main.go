package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/alexptdev/chat-server-api/internal/config"
	"github.com/alexptdev/chat-server-api/internal/config/env"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"

	desc "github.com/alexptdev/chat-server-api/pkg/chat_v1"
)

var configPath string

type chatServer struct {
	desc.UnimplementedChatV1Server
	conPool *pgxpool.Pool
}

func init() {
	flag.StringVar(
		&configPath,
		"config-path",
		".env",
		"path to config file",
	)
}

func (s *chatServer) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {

	insertBuilder := sq.Insert("chats").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_name", "chat_description", "chat_author_id").
		Values(req.Name, req.Description, req.AuthorId).
		Suffix("RETURNING chat_id")

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v \n", err)
		return nil, err
	}

	var chatId int64
	err = s.conPool.QueryRow(ctx, query, args...).Scan(&chatId)
	if err != nil {
		log.Printf("failed to create chat: %v \n", err)
		return nil, err
	}

	return &desc.CreateResponse{Id: chatId}, nil
}

func (s *chatServer) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {

	deleteBuilder := sq.Delete("chats").
		Where(sq.Eq{"chat_id": req.Id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v \n", err)
		return nil, err
	}

	_, err = s.conPool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to delete chat: %v \n", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *chatServer) AddUser(ctx context.Context, req *desc.AddUserRequest) (*emptypb.Empty, error) {

	insertBuilder := sq.Insert("chat_users").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_users_chat_id", "chat_users_user_id").
		Values(req.ChatId, req.UserId)

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v \n", err)
		return nil, err
	}

	_, err = s.conPool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to add user to chat: %v \n", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *chatServer) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {

	insertBuilder := sq.Insert("chat_messages").
		PlaceholderFormat(sq.Dollar).
		Columns("message_chat_id", "message_from", "message_text").
		Values(req.ChatId, req.From, req.Text)

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v \n", err)
		return nil, err
	}

	_, err = s.conPool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to add user to chat: %v \n", err)
	}

	return &emptypb.Empty{}, nil
}

func main() {

	ctx := context.Background()

	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config %v", err)
	}

	grpcConfig, err := env.NewGrpcConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPgConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	listen, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("Failed to listen port: %v", err)
	}

	conPool, err := pgxpool.Connect(ctx, pgConfig.Dsn())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conPool.Close()

	server := grpc.NewServer()
	reflection.Register(server)
	desc.RegisterChatV1Server(server, &chatServer{
		conPool: conPool,
	})

	log.Printf("Server listening at %v", listen.Addr())
	if err = server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
