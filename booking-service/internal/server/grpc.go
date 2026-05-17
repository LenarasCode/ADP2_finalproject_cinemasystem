package server

import (
	"github.com/cinema-system/booking-service/internal/service"
	pb "github.com/cinema-system/booking-service/proto/booking"
)

func NewBookingServer(svc *service.BookingService) pb.BookingServiceServer {
	return svc
}
