package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/cinema-system/api-gateway/internal/handlers"
	pb "github.com/cinema-system/api-gateway/proto/gateway"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50050"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCinemaGatewayServer(s, handlers.NewGatewayServer())
	reflection.Register(s)

	log.Printf("API Gateway listening on :%s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
