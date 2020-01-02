// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"net"

	pb "bitbucket.org/sothy5/n4-go/client/pkg/api"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

//protoc --proto_path=api/proto  --go_out=plugins=grpc:pkg/api/ session_establishment.proto
// server is used to implement SessionEstablishmentServiceServer
type server struct {
	pb.UnimplementedSessionEstablishmentServiceServer
}

// Create implements SessionEstablishmentServiceServer
func (s *server) Create(ctx context.Context, in *pb.SessionEstablishmentRequest) (*pb.SessionEstablishmentResponse, error) {
	log.Printf("Received supi: %v", in.GetSupi())
	cnTunnelInfo := pb.CNTunnelInfo{
		TunnelId:    101,
		Ipv4Address: "192.168.1.104",
	}
	return &pb.SessionEstablishmentResponse{Supi: in.GetSupi(),
		CnTunnelInfo: &cnTunnelInfo,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSessionEstablishmentServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
