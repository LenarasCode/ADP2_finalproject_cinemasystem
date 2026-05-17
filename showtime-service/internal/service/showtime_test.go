package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cinema-system/showtime-service/internal/repository"
	"github.com/cinema-system/showtime-service/internal/service"
	pb "github.com/cinema-system/showtime-service/proto/showtime"
)

// mockRepo реализует service.ShowtimeRepository
type mockRepo struct {
	getShowtimeFn          func(string) (*repository.Showtime, error)
	listShowtimesByMovieFn func(string) ([]repository.Showtime, error)
	getSeatsFn             func(string) ([]repository.Seat, error)
	createShowtimeFn       func(*repository.Showtime) (*repository.Showtime, error)
	updateShowtimeFn       func(*repository.Showtime) (*repository.Showtime, error)
	deleteShowtimeFn       func(string) error
}

func (m *mockRepo) GetShowtime(id string) (*repository.Showtime, error) {
	return m.getShowtimeFn(id)
}
func (m *mockRepo) ListShowtimesByMovie(movieID string) ([]repository.Showtime, error) {
	return m.listShowtimesByMovieFn(movieID)
}
func (m *mockRepo) GetSeats(showtimeID string) ([]repository.Seat, error) {
	return m.getSeatsFn(showtimeID)
}
func (m *mockRepo) CreateShowtime(s *repository.Showtime) (*repository.Showtime, error) {
	return m.createShowtimeFn(s)
}
func (m *mockRepo) UpdateShowtime(s *repository.Showtime) (*repository.Showtime, error) {
	return m.updateShowtimeFn(s)
}
func (m *mockRepo) DeleteShowtime(id string) error {
	return m.deleteShowtimeFn(id)
}

// mockCache реализует service.ShowtimeCache
type mockCache struct{}

func (m *mockCache) GetShowtime(ctx context.Context, id string) (*repository.Showtime, error) {
	return nil, errors.New("cache miss")
}
func (m *mockCache) SetShowtime(ctx context.Context, s *repository.Showtime) error { return nil }
func (m *mockCache) DeleteShowtime(ctx context.Context, id string) error            { return nil }

func TestGetShowtime_Success(t *testing.T) {
	repo := &mockRepo{
		getShowtimeFn: func(id string) (*repository.Showtime, error) {
			return &repository.Showtime{
				ID: "123", MovieID: "m1", Hall: "Hall 1", StartTime: "2025-01-01T10:00:00Z", AvailableSeats: 50,
			}, nil
		},
	}
	svc := service.NewShowtimeService(repo, &mockCache{})

	resp, err := svc.GetShowtime(context.Background(), &pb.GetShowtimeRequest{Id: "123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Id != "123" || resp.AvailableSeats != 50 {
		t.Errorf("got %+v", resp)
	}
}

func TestGetShowtime_NotFound(t *testing.T) {
	repo := &mockRepo{
		getShowtimeFn: func(id string) (*repository.Showtime, error) {
			return nil, errors.New("not found")
		},
	}
	svc := service.NewShowtimeService(repo, &mockCache{})

	_, err := svc.GetShowtime(context.Background(), &pb.GetShowtimeRequest{Id: "bad"})
	if err == nil {
		t.Error("expected error, got nil")
	}
}
