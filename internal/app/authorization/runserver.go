package authServer

import (
	"2021_2_MAMBa/internal/app/config"
	"2021_2_MAMBa/internal/pkg/sessions"
	grpcSessionServer "2021_2_MAMBa/internal/pkg/sessions/delivery/grpc"
	"2021_2_MAMBa/internal/pkg/utils/log"
	"google.golang.org/grpc"
	"net"
)

func RunServer(configPath string) {
	cfg := config.ParseAuth(configPath)
	lis, err := net.Listen("tcp", "localhost:"+cfg.AuthPort)
	if err != nil {
		log.Warn("failed to listen:"+ err.Error())
	}
	s := grpc.NewServer()
	grpcSessionServer.RegisterSessionRPCServer(s, sessions.NewSessionManager(cfg.Secure))
	log.Info("server listening at"+ lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		log.Warn("failed to listen:"+ err.Error())
	}
}
