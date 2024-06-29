package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"

	protos "farseer/protos"
	"farseer/time"

	"github.com/charmbracelet/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const PORT = 2285

type hubRPCServer struct {
	// utils
	netw Network
	ll   log.Logger

	protos.UnimplementedHubServiceServer
	rpcServer map[string][]*protos.HubServiceServer
}

func (s *hubRPCServer) SubmitMessage(ctx context.Context, message *protos.Message) (*protos.Message, error) {
	peerIdEncoded, err := s.netw.self.Marshal()
	if err != nil {
		return &protos.Message{}, err
	}

	contactInfoTime, err := time.GetFarcasterTime()
	if err != nil {
		contactInfoTime = 0
		s.netw.logger.Error("Couldn't get Farcaster time for the message!")
	}

	msg := protos.GossipMessage{
		Content: &protos.GossipMessage_Message{
			Message: message,
		},
		Topics:    []string{s.netw.topic.String()},
		PeerId:    peerIdEncoded,
		Version:   protos.GossipVersion_GOSSIP_VERSION_V1_1,
		Timestamp: uint32(contactInfoTime),
	}

	s.netw.Publish(&msg)

	return message, nil
}

func (s *hubRPCServer) ValidateMessage(ctx context.Context, message *protos.Message) (*protos.ValidationResponse, error) {
	return &protos.ValidationResponse{Valid: true}, nil
}

func newServer(netw Network, ll log.Logger) *hubRPCServer {
	s := &hubRPCServer{
		netw:      netw,
		ll:        ll,
		rpcServer: make(map[string][]*protos.HubServiceServer),
	}
	return s
}

func Start(wg *sync.WaitGroup, stopCh <-chan struct{}, netw Network) {
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
	protos.RegisterHubServiceServer(grpcServer, newServer(netw, *ll))
	reflection.Register(grpcServer)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			ll.Fatal("Failed to serve:", err)
		}
	}()

	<-stopCh

	grpcServer.GracefulStop()
	ll.Info("Graceful shutdown was successful!")
}
