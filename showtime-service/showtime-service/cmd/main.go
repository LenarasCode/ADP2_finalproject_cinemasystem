package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"github.com/cinema-system/showtime-service/internal/cache"
	"github.com/cinema-system/showtime-service/internal/repository"
	"github.com/cinema-system/showtime-service/internal/server"
	"github.com/cinema-system/showtime-service/internal/service"
	pb "github.com/cinema-system/showtime-service/proto/showtime"
	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://showtime:showtime@localhost:5433/showtimedb?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})

	repo := repository.NewShowtimeRepo(db)
	cch := cache.NewShowtimeCache(rdb)
	svc := service.NewShowtimeService(repo, cch)

	grpcServer := grpc.NewServer()
	pb.RegisterShowtimeServiceServer(grpcServer, server.NewShowtimeServer(svc))

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50052"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Showtime service listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
