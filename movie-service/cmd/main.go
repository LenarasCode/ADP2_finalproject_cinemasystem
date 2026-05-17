package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/cinema-system/movie-service/internal/cache"
	"github.com/cinema-system/movie-service/internal/repository"
	"github.com/cinema-system/movie-service/internal/server"
	"github.com/cinema-system/movie-service/internal/service"
	pb "github.com/cinema-system/movie-service/proto/movie"
	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://movie:movie@localhost:5434/moviedb?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer db.Close()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})

	repo := repository.NewMovieRepo(db)
	cch := cache.NewMovieCache(rdb)
	svc := service.NewMovieService(repo, cch)

	s := grpc.NewServer()
	pb.RegisterMovieServiceServer(s, server.NewMovieServer(svc))
	reflection.Register(s)

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("Movie service listening on :%s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
