package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/cinema-system/booking-service/internal/publisher"
	"github.com/cinema-system/booking-service/internal/repository"
	"github.com/cinema-system/booking-service/internal/server"
	"github.com/cinema-system/booking-service/internal/service"
	pb "github.com/cinema-system/booking-service/proto/booking"
	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://booking:booking@localhost:5432/bookingdb?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer db.Close()

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	pub, err := publisher.NewNATSPublisher(natsURL)
	if err != nil {
		log.Printf("NATS not available, continuing without it: %v", err)
		pub = nil
	}
	defer pub.Close()

	repo := repository.NewBookingRepo(db)
	svc := service.NewBookingService(repo, db, pub)

	s := grpc.NewServer()
	pb.RegisterBookingServiceServer(s, server.NewBookingServer(svc))
	reflection.Register(s)

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50053"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("Booking service listening on :%s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
