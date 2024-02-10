package main

import (
	"context"
	"fmt"
	desc "github.com/alexptdev/chat-server-api/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

const grpcPort = 4001

type chatServer struct {
	desc.UnimplementedChatV1Server
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

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
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
