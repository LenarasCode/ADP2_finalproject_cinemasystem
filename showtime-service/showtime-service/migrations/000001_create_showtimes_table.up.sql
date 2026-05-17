CREATE TABLE showtimes (
    id UUID PRIMARY KEY,
    movie_id VARCHAR(255) NOT NULL,
    hall VARCHAR(100) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    available_seats INT NOT NULL
);
