package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/alexptdev/chat-server-api/internal/config"
	"github.com/alexptdev/chat-server-api/internal/config/env"
	desc "github.com/alexptdev/chat-server-api/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

var configPath string

type chatServer struct {
	desc.UnimplementedChatV1Server
}

func init() {
	flag.StringVar(
		&configPath,
		"config-path",
		".env",
		"path to config file",
	)
}

func (s *chatServer) Create(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {

	for _, item := range req.UserNames {
		fmt.Println(item)
	}

	return &desc.CreateResponse{Id: 1}, nil
}

func (s *chatServer) Delete(_ context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	fmt.Println(req.Id)
	return &emptypb.Empty{}, nil
}

func (s *chatServer) SendMessage(_ context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	fmt.Println(req.From)
	fmt.Println(req.Text)
	fmt.Println(req.Timestamp.AsTime())
	return &emptypb.Empty{}, nil
}

func main() {

	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config %v", err)
	}

	grpcConfig, err := env.NewGrpcConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	_, err = env.NewPgConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	listen, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("Failed to listen port: %v", err)
	}

	server := grpc.NewServer()
	reflection.Register(server)
	desc.RegisterChatV1Server(server, &chatServer{})

	log.Printf("Server listening at %v", listen.Addr())
	if err = server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
