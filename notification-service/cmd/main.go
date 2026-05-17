package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/cinema-system/notification-service/internal/email"
	"github.com/cinema-system/notification-service/internal/nats"
	"github.com/cinema-system/notification-service/internal/repository"
	"github.com/cinema-system/notification-service/internal/server"
	"github.com/cinema-system/notification-service/internal/service"
	pb "github.com/cinema-system/notification-service/proto/notification"
	_ "github.com/lib/pq"
)

func main() {
	// PostgreSQL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://notification:notification@localhost:5435/notificationdb?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer db.Close()

	// NATS
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("nats connect: %v", err)
	}
	defer nc.Close()

	// Email sender
	sender := email.NewSender()

	repo := repository.NewEmailLogRepo(db)
	svc := service.NewNotificationService(repo, sender)

	// Запускаем NATS подписчик
	go nats_subscribe.Subscribe(nc, svc)

	// gRPC сервер
	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, server.NewNotificationServer(svc))
	reflection.Register(grpcServer)

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50054"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("Notification service listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
