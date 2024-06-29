package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"

	protos "farseer/protos"

	"github.com/charmbracelet/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const PORT = 2285


type hubRPCServer struct {
	protos.UnimplementedHubServiceServer
	rpcServer map[string][]*protos.HubServiceServer
}

func (s *hubRPCServer) SubmitMessage(ctx context.Context, message *protos.Message) (*protos.Message, error) {
	return message, nil
}

func (s *hubRPCServer) ValidateMessage(ctx context.Context, message *protos.Message) (*protos.ValidationResponse, error) {
	return &protos.ValidationResponse{Valid: true}, nil
}

func newServer() *hubRPCServer {
	s := &hubRPCServer{
		rpcServer: make(map[string][]*protos.HubServiceServer),
	}
	return s
}

func Start(wg *sync.WaitGroup, stopCh <-chan struct{}) {
	defer wg.Done()

	ll := log.New(os.Stderr)
	ll.SetPrefix("grpc")

	// todo: read from config
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", PORT))
	if err != nil {
		ll.Fatal("Can't start the listnener! |", "Err", err)
	}

	ll.Info("Started the server!")

	grpcServer := grpc.NewServer(grpc.EmptyServerOption{})
	protos.RegisterHubServiceServer(grpcServer, newServer())
	reflection.Register(grpcServer)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			ll.Fatal("Failed to serve:", err)
		}	
	}()

	<- stopCh

	grpcServer.GracefulStop()
	ll.Info("Graceful shutdown was successful!")
}