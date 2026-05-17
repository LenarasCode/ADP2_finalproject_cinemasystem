package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cinema-system/booking-service/internal/repository"
	"github.com/cinema-system/booking-service/internal/service"
	pb "github.com/cinema-system/booking-service/proto/booking"
)

type mockRepo struct {
	createBookingFn func(tx interface{}, userID, showtimeID string, seatIDs []string) (*repository.Booking, error)
	getBookingFn    func(id string) (*repository.Booking, []string, error)
	cancelBookingFn func(tx interface{}, bookingID string) error
}
func (m *mockRepo) CreateBooking(tx interface{}, userID, showtimeID string, seatIDs []string) (*repository.Booking, error) {
	return m.createBookingFn(tx, userID, showtimeID, seatIDs)
}
func (m *mockRepo) GetBooking(id string) (*repository.Booking, []string, error) {
	return m.getBookingFn(id)
}
func (m *mockRepo) CancelBooking(tx interface{}, bookingID string) error {
	return m.cancelBookingFn(tx, bookingID)
}

func TestCreateBooking_Success(t *testing.T) {
	repo := &mockRepo{
		createBookingFn: func(tx interface{}, userID, showtimeID string, seatIDs []string) (*repository.Booking, error) {
			return &repository.Booking{ID: "b1", UserID: userID, ShowtimeID: showtimeID, Status: "CONFIRMED"}, nil
		},
	}
	svc := service.NewBookingService(repo, nil, nil)
	resp, err := svc.CreateBooking(context.Background(), &pb.CreateBookingRequest{
		UserId: "u1", ShowtimeId: "s1", SeatIds: []string{"seat1"},
	})
	if err != nil || resp.Status != "CONFIRMED" {
		t.Errorf("expected CONFIRMED, got %v, err %v", resp, err)
	}
}

func TestCancelBooking_Success(t *testing.T) {
	repo := &mockRepo{
		cancelBookingFn: func(tx interface{}, bookingID string) error { return nil },
	}
	svc := service.NewBookingService(repo, nil, nil)
	_, err := svc.CancelBooking(context.Background(), &pb.CancelBookingRequest{Id: "b1"})
	if err != nil {
		t.Error(err)
	}
}
