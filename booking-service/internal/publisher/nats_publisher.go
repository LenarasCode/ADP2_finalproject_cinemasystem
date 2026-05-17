package publisher

import (
	"log"
	"github.com/nats-io/nats.go"
)

type NATSPublisher struct {
	conn *nats.Conn
}

func NewNATSPublisher(url string) (*NATSPublisher, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NATSPublisher{conn: nc}, nil
}

func (p *NATSPublisher) PublishBookingCreated(bookingID string) {
	if err := p.conn.Publish("booking.created", []byte(bookingID)); err != nil {
		log.Printf("Failed to publish NATS event: %v", err)
	} else {
		log.Printf("Published booking.created for booking %s", bookingID)
	}
}

func (p *NATSPublisher) Close() {
	p.conn.Close()
}
