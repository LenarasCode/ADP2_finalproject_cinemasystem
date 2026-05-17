package nats

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/cinema-system/notification-service/internal/service"
	pb "github.com/cinema-system/notification-service/proto/notification"
)

func Subscribe(nc *nats.Conn, svc *service.NotificationService) {
	_, err := nc.Subscribe("booking.created", func(msg *nats.Msg) {
		// В реальности здесь парсим booking_id и получаем email.
		// Пока для теста отправляем на фиксированный адрес.
		bookingID := string(msg.Data)
		log.Printf("Received NATS event for booking %s", bookingID)
		_, err := svc.SendBookingConfirmation(nil, &pb.SendBookingConfirmationRequest{
			BookingId:     bookingID,
			RecipientEmail: "customer@example.com", // заглушка
		})
		if err != nil {
			log.Printf("Failed to send confirmation: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to NATS: %v", err)
	}
	log.Println("NATS subscribed to booking.created")
}
