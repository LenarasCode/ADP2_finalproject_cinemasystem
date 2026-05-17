package server

import (
	"github.com/cinema-system/notification-service/internal/service"
	pb "github.com/cinema-system/notification-service/proto/notification"
)

func NewNotificationServer(svc *service.NotificationService) pb.NotificationServiceServer {
	return svc
}
