package service

import (
	"context"
	"log"

	"github.com/cinema-system/showtime-service/internal/repository"
	pb "github.com/cinema-system/showtime-service/proto/showtime"
)

type ShowtimeRepository interface {
	GetShowtime(id string) (*repository.Showtime, error)
	ListShowtimesByMovie(movieID string) ([]repository.Showtime, error)
	GetSeats(showtimeID string) ([]repository.Seat, error)
	CreateShowtime(s *repository.Showtime) (*repository.Showtime, error)
	UpdateShowtime(s *repository.Showtime) (*repository.Showtime, error)
	DeleteShowtime(id string) error
}

type ShowtimeCache interface {
	GetShowtime(ctx context.Context, id string) (*repository.Showtime, error)
	SetShowtime(ctx context.Context, s *repository.Showtime) error
	DeleteShowtime(ctx context.Context, id string) error
}

type ShowtimeService struct {
	repo  ShowtimeRepository
	cache ShowtimeCache
	pb.UnimplementedShowtimeServiceServer
}

func NewShowtimeService(repo ShowtimeRepository, cache ShowtimeCache) *ShowtimeService {
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

func (s *ShowtimeService) CreateShowtime(ctx context.Context, req *pb.CreateShowtimeRequest) (*pb.Showtime, error) {
	st := &repository.Showtime{
		MovieID:        req.MovieId,
		Hall:           req.Hall,
		StartTime:      req.StartTime,
		AvailableSeats: int(req.AvailableSeats),
	}
	created, err := s.repo.CreateShowtime(st)
	if err != nil {
		return nil, err
	}
	return repoToProto(created), nil
}

func (s *ShowtimeService) UpdateShowtime(ctx context.Context, req *pb.UpdateShowtimeRequest) (*pb.Showtime, error) {
	st := &repository.Showtime{
		ID:             req.Id,
		MovieID:        req.MovieId,
		Hall:           req.Hall,
		StartTime:      req.StartTime,
		AvailableSeats: int(req.AvailableSeats),
	}
	updated, err := s.repo.UpdateShowtime(st)
	if err != nil {
		return nil, err
	}
	go func() {
		if err := s.cache.DeleteShowtime(context.Background(), req.Id); err != nil {
			log.Printf("cache delete error: %v", err)
		}
	}()
	return repoToProto(updated), nil
}

func (s *ShowtimeService) DeleteShowtime(ctx context.Context, req *pb.DeleteShowtimeRequest) (*pb.DeleteShowtimeResponse, error) {
	err := s.repo.DeleteShowtime(req.Id)
	if err != nil {
		return nil, err
	}
	go func() {
		if err := s.cache.DeleteShowtime(context.Background(), req.Id); err != nil {
			log.Printf("cache delete error: %v", err)
		}
	}()
	return &pb.DeleteShowtimeResponse{Success: true}, nil
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
