package authServer

import (
	"2021_2_MAMBa/internal/pkg/sessions"
	grpcSessionServer "2021_2_MAMBa/internal/pkg/sessions/delivery/grpc"
	"google.golang.org/grpc"
	"log"
	"net"
)

func RunServer(addr string) {
	lis, err := net.Listen("tcp", "localhost:"+addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	grpcSessionServer.RegisterSessionRPCServer(s, sessions.NewSessionManager())
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
