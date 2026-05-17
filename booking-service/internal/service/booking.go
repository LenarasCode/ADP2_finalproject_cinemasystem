package service

import (
	"context"
	"database/sql"

	"github.com/cinema-system/booking-service/internal/publisher"
	"github.com/cinema-system/booking-service/internal/repository"
	pb "github.com/cinema-system/booking-service/proto/booking"
)

type BookingService struct {
	repo      *repository.BookingRepo
	db        *sql.DB
	publisher *publisher.NATSPublisher
	pb.UnimplementedBookingServiceServer
}

func NewBookingService(repo *repository.BookingRepo, db *sql.DB, pub *publisher.NATSPublisher) *BookingService {
	return &BookingService{repo: repo, db: db, publisher: pub}
}

func (s *BookingService) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.Booking, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	booking, err := s.repo.CreateBooking(tx, req.UserId, req.ShowtimeId, req.SeatIds)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Публикуем событие в NATS
	if s.publisher != nil {
		s.publisher.PublishBookingCreated(booking.ID)
	}

	return &pb.Booking{
		Id:         booking.ID,
		UserId:     booking.UserID,
		ShowtimeId: booking.ShowtimeID,
		SeatIds:    req.SeatIds,
		Status:     booking.Status,
	}, nil
}

func (s *BookingService) GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.Booking, error) {
	booking, seats, err := s.repo.GetBooking(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Booking{
		Id:         booking.ID,
		UserId:     booking.UserID,
		ShowtimeId: booking.ShowtimeID,
		SeatIds:    seats,
		Status:     booking.Status,
	}, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, req *pb.CancelBookingRequest) (*pb.CancelBookingResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := s.repo.CancelBooking(tx, req.Id); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &pb.CancelBookingResponse{Success: true}, nil
}
