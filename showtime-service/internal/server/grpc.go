package server

import (
	"github.com/cinema-system/showtime-service/internal/service"
	pb "github.com/cinema-system/showtime-service/proto/showtime"
)

func NewShowtimeServer(svc *service.ShowtimeService) pb.ShowtimeServiceServer {
	return svc
}
