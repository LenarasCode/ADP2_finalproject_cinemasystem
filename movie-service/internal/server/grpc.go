package server

import (
	"github.com/cinema-system/movie-service/internal/service"
	pb "github.com/cinema-system/movie-service/proto/movie"
)

func NewMovieServer(svc *service.MovieService) pb.MovieServiceServer {
	return svc
}
