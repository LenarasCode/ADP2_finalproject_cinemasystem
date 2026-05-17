package service

import (
	"context"
	"log"
	"time"

	"github.com/cinema-system/notification-service/internal/email"
	"github.com/cinema-system/notification-service/internal/repository"
	pb "github.com/cinema-system/notification-service/proto/notification"
	"github.com/google/uuid"
)

type NotificationService struct {
	repo   *repository.EmailLogRepo
	sender *email.Sender
	pb.UnimplementedNotificationServiceServer
}

func NewNotificationService(repo *repository.EmailLogRepo, sender *email.Sender) *NotificationService {
	return &NotificationService{repo: repo, sender: sender}
}

func (s *NotificationService) SendBookingConfirmation(ctx context.Context, req *pb.SendBookingConfirmationRequest) (*pb.SendBookingConfirmationResponse, error) {
	subject := "Booking Confirmation"
	body := "Your booking " + req.BookingId + " is confirmed."
	err := s.sender.Send(req.RecipientEmail, subject, body)
	status := "OK"
	if err != nil {
		status = "FAILED: " + err.Error()
		log.Printf("Failed to send email: %v", err)
	}
	// log to DB
	logEntry := &repository.EmailLog{
		ID:        uuid.New().String(),
		Recipient: req.RecipientEmail,
		Subject:   subject,
		Body:      body,
		SentAt:    time.Now(),
		Status:    status,
	}
	if dbErr := s.repo.Insert(logEntry); dbErr != nil {
		log.Printf("Failed to insert email log: %v", dbErr)
	}
	return &pb.SendBookingConfirmationResponse{Success: err == nil}, err
}

func (s *NotificationService) GetEmailLogs(ctx context.Context, req *pb.GetEmailLogsRequest) (*pb.GetEmailLogsResponse, error) {
	logs, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	var pbLogs []*pb.EmailLog
	for _, l := range logs {
		pbLogs = append(pbLogs, &pb.EmailLog{
			Id:        l.ID,
			Recipient: l.Recipient,
			Subject:   l.Subject,
			Body:      l.Body,
			SentAt:    l.SentAt.Format(time.RFC3339),
			Status:    l.Status,
		})
	}
	return &pb.GetEmailLogsResponse{Logs: pbLogs}, nil
}

func (s *NotificationService) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Ok: true}, nil
}
