package collectionsServer

import (
	"2021_2_MAMBa/internal/app/config"
	grpcCollectionServer "2021_2_MAMBa/internal/pkg/collections/delivery/grpc"
	collectionsRepository "2021_2_MAMBa/internal/pkg/collections/repository"
	collectionsUsecase "2021_2_MAMBa/internal/pkg/collections/usecase"
	"2021_2_MAMBa/internal/pkg/database"
	"google.golang.org/grpc"
	"log"
	"net"
)

func RunServer(configPath string) {
	cfg := config.ParseMain(configPath)
	db := database.Connect(cfg.Db)
	defer db.Disconnect()
	collectionsRepo := collectionsRepository.NewCollectionsRepository(db)
	colUsecase := collectionsUsecase.NewCollectionsUsecase(collectionsRepo)

	lis, err := net.Listen("tcp", "localhost:"+cfg.CollectPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	grpcCollectionServer.RegisterCollectionsRPCServer(s, colUsecase)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
