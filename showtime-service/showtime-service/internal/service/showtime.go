package service

import (
	"context"
	"log"

	"github.com/cinema-system/showtime-service/internal/cache"
	"github.com/cinema-system/showtime-service/internal/repository"
	pb "github.com/cinema-system/showtime-service/proto/showtime"
)

type ShowtimeService struct {
	repo  *repository.ShowtimeRepo
	cache *cache.ShowtimeCache
	pb.UnimplementedShowtimeServiceServer
}

func NewShowtimeService(repo *repository.ShowtimeRepo, cache *cache.ShowtimeCache) *ShowtimeService {
	return &ShowtimeService{repo: repo, cache: cache}
}

func (s *ShowtimeService) GetShowtime(ctx context.Context, req *pb.GetShowtimeRequest) (*pb.Showtime, error) {
	cached, err := s.cache.GetShowtime(ctx, req.Id)
	if err == nil {
		return repoToProto(cached), nil
	}

	st, err := s.repo.GetShowtime(req.Id)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := s.cache.SetShowtime(context.Background(), st); err != nil {
			log.Printf("cache set error: %v", err)
		}
	}()

	return repoToProto(st), nil
}

func (s *ShowtimeService) ListShowtimesByMovie(ctx context.Context, req *pb.ListShowtimesByMovieRequest) (*pb.ListShowtimesResponse, error) {
	showtimes, err := s.repo.ListShowtimesByMovie(req.MovieId)
	if err != nil {
		return nil, err
	}
	var pbList []*pb.Showtime
	for i := range showtimes {
		pbList = append(pbList, repoToProto(&showtimes[i]))
	}
	return &pb.ListShowtimesResponse{Showtimes: pbList}, nil
}

func (s *ShowtimeService) GetSeatAvailability(ctx context.Context, req *pb.GetSeatAvailabilityRequest) (*pb.SeatAvailabilityResponse, error) {
	seats, err := s.repo.GetSeats(req.ShowtimeId)
	if err != nil {
		return nil, err
	}
	var pbSeats []*pb.Seat
	for _, seat := range seats {
		pbSeats = append(pbSeats, &pb.Seat{
			Id:          seat.ID,
			Row:         seat.Row,
			Number:      int32(seat.Number),
			IsAvailable: seat.IsAvailable,
		})
	}
	return &pb.SeatAvailabilityResponse{Seats: pbSeats}, nil
}

func repoToProto(s *repository.Showtime) *pb.Showtime {
	return &pb.Showtime{
		Id:             s.ID,
		MovieId:        s.MovieID,
		Hall:           s.Hall,
		StartTime:      s.StartTime,
		AvailableSeats: int32(s.AvailableSeats),
	}
}
