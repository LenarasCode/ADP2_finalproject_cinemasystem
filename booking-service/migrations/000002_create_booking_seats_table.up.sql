CREATE TABLE booking_seats (
    booking_id UUID NOT NULL REFERENCES bookings(id),
    seat_id UUID NOT NULL,
    PRIMARY KEY (booking_id, seat_id)
);
