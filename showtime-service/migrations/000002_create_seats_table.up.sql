CREATE TABLE seats (
    id UUID PRIMARY KEY,
    showtime_id UUID NOT NULL REFERENCES showtimes(id),
    row VARCHAR(5) NOT NULL,
    number INT NOT NULL,
    is_available BOOLEAN DEFAULT true
);
