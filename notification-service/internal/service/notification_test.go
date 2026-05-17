package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cinema-system/notification-service/internal/repository"
	"github.com/cinema-system/notification-service/internal/service"
	pb "github.com/cinema-system/notification-service/proto/notification"
)

type mockRepo struct {
	insertFn func(log *repository.EmailLog) error
	getAllFn func() ([]repository.EmailLog, error)
}
func (m *mockRepo) Insert(log *repository.EmailLog) error { return m.insertFn(log) }
func (m *mockRepo) GetAll() ([]repository.EmailLog, error) { return m.getAllFn() }

type mockSender struct {
	sendFn func(to, subject, body string) error
}
func (m *mockSender) Send(to, subject, body string) error { return m.sendFn(to, subject, body) }

func TestSendBookingConfirmation_Success(t *testing.T) {
	repo := &mockRepo{insertFn: func(log *repository.EmailLog) error { return nil }}
	sender := &mockSender{sendFn: func(to, subject, body string) error { return nil }}
	svc := service.NewNotificationService(repo, sender)

	resp, err := svc.SendBookingConfirmation(context.Background(), &pb.SendBookingConfirmationRequest{
		BookingId:     "123",
		RecipientEmail: "a@b.com",
	})
	if err != nil || !resp.Success {
		t.Errorf("expected success, got %v, %v", resp, err)
	}
}

func TestSendBookingConfirmation_SendFail(t *testing.T) {
	repo := &mockRepo{insertFn: func(log *repository.EmailLog) error { return nil }}
	sender := &mockSender{sendFn: func(to, subject, body string) error { return errors.New("smtp error") }}
	svc := service.NewNotificationService(repo, sender)

	resp, err := svc.SendBookingConfirmation(context.Background(), &pb.SendBookingConfirmationRequest{
		BookingId:     "123",
		RecipientEmail: "a@b.com",
	})
	if err == nil || resp.Success {
		t.Errorf("expected failure, got %v, %v", resp, err)
	}
}
