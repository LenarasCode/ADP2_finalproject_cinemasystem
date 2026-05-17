CREATE TABLE movies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    genre VARCHAR(100) NOT NULL,
    duration_min INT NOT NULL,
    poster_url TEXT
);
