package handlers

import (
	"context"
	"log"

	pb "github.com/cinema-system/api-gateway/proto/gateway"
)

type GatewayServer struct {
	pb.UnimplementedCinemaGatewayServer
}

func NewGatewayServer() *GatewayServer {
	return &GatewayServer{}
}

// Метод Ленары – GetShowtimeWithMovie
func (s *GatewayServer) GetShowtimeWithMovie(ctx context.Context, req *pb.GetShowtimeWithMovieRequest) (*pb.ShowtimeWithMovie, error) {
	// Здесь должен быть вызов Showtime Service и Movie Service, но пока возвращаем заглушку
	log.Printf("Ленара: запрос информации о сеансе %s", req.ShowtimeId)
	return &pb.ShowtimeWithMovie{
		ShowtimeId: req.ShowtimeId,
		MovieTitle: "The Matrix",
		Hall:       "Hall 1",
		StartTime:  "2026-06-01T12:00:00Z",
	}, nil
}

// Метод Турара – SendBookingNotification
func (s *GatewayServer) SendBookingNotification(ctx context.Context, req *pb.SendBookingNotificationRequest) (*pb.SendBookingNotificationResponse, error) {
	log.Printf("Турар: отправка уведомления для бронирования %s на %s", req.BookingId, req.RecipientEmail)
	// Вызов Notification Service
	return &pb.SendBookingNotificationResponse{Success: true}, nil
}

// Метод Алмаса – GetBookingDetails
func (s *GatewayServer) GetBookingDetails(ctx context.Context, req *pb.GetBookingDetailsRequest) (*pb.BookingDetails, error) {
	log.Printf("Алмас: запрос деталей бронирования %s", req.BookingId)
	// Вызов Booking Service + Showtime + Movie
	return &pb.BookingDetails{
		BookingId:  req.BookingId,
		UserId:     "u1",
		MovieTitle: "Inception",
		Hall:       "Hall 3",
		StartTime:  "2026-06-02T14:00:00Z",
		Seats:      []string{"A1", "A2"},
		Status:     "CONFIRMED",
	}, nil
}
