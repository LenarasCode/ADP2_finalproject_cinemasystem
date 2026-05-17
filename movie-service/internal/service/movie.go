package service

import (
	"context"
	"log"

	"github.com/cinema-system/movie-service/internal/cache"
	"github.com/cinema-system/movie-service/internal/repository"
	pb "github.com/cinema-system/movie-service/proto/movie"
)

type MovieService struct {
	repo  *repository.MovieRepo
	cache *cache.MovieCache
	pb.UnimplementedMovieServiceServer
}

func NewMovieService(repo *repository.MovieRepo, cache *cache.MovieCache) *MovieService {
	return &MovieService{repo: repo, cache: cache}
}

func (s *MovieService) GetMovie(ctx context.Context, req *pb.GetMovieRequest) (*pb.Movie, error) {
	cached, err := s.cache.GetMovie(ctx, req.Id)
	if err == nil {
		return movieToProto(cached), nil
	}
	m, err := s.repo.GetMovie(req.Id)
	if err != nil {
		return nil, err
	}
	go func() {
		if err := s.cache.SetMovie(context.Background(), m); err != nil {
			log.Printf("cache set error: %v", err)
		}
	}()
	return movieToProto(m), nil
}

func (s *MovieService) ListMovies(ctx context.Context, req *pb.ListMoviesRequest) (*pb.ListMoviesResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)
	if page < 1 { page = 1 }
	if limit < 1 || limit > 50 { limit = 10 }
	movies, total, err := s.repo.ListMovies(page, limit)
	if err != nil {
		return nil, err
	}
	var pbMovies []*pb.Movie
	for i := range movies {
		pbMovies = append(pbMovies, movieToProto(&movies[i]))
	}
	return &pb.ListMoviesResponse{Movies: pbMovies, Total: int32(total)}, nil
}

func (s *MovieService) SearchMovies(ctx context.Context, req *pb.SearchMoviesRequest) (*pb.ListMoviesResponse, error) {
	movies, err := s.repo.SearchMovies(req.Query)
	if err != nil {
		return nil, err
	}
	var pbMovies []*pb.Movie
	for i := range movies {
		pbMovies = append(pbMovies, movieToProto(&movies[i]))
	}
	return &pb.ListMoviesResponse{Movies: pbMovies, Total: int32(len(pbMovies))}, nil
}

func movieToProto(m *repository.Movie) *pb.Movie {
	return &pb.Movie{
		Id:          m.ID,
		Title:       m.Title,
		Genre:       m.Genre,
		DurationMin: int32(m.DurationMin),
		PosterUrl:   m.PosterURL,
	}
}
