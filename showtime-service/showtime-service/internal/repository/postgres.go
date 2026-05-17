package repository

import (
	"database/sql"
	"fmt"
)

type ShowtimeRepo struct {
	db *sql.DB
}

func NewShowtimeRepo(db *sql.DB) *ShowtimeRepo {
	return &ShowtimeRepo{db: db}
}

type Showtime struct {
	ID             string
	MovieID        string
	Hall           string
	StartTime      string
	AvailableSeats int
}

type Seat struct {
	ID          string
	Row         string
	Number      int
	IsAvailable bool
	ShowtimeID  string
}

func (r *ShowtimeRepo) GetShowtime(id string) (*Showtime, error) {
	row := r.db.QueryRow(
		"SELECT id, movie_id, hall, start_time, available_seats FROM showtimes WHERE id=$1", id,
	)
	var s Showtime
	if err := row.Scan(&s.ID, &s.MovieID, &s.Hall, &s.StartTime, &s.AvailableSeats); err != nil {
		return nil, fmt.Errorf("showtime not found: %w", err)
	}
	return &s, nil
}

func (r *ShowtimeRepo) ListShowtimesByMovie(movieID string) ([]Showtime, error) {
	rows, err := r.db.Query(
		"SELECT id, movie_id, hall, start_time, available_seats FROM showtimes WHERE movie_id=$1 ORDER BY start_time",
		movieID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var showtimes []Showtime
	for rows.Next() {
		var s Showtime
		if err := rows.Scan(&s.ID, &s.MovieID, &s.Hall, &s.StartTime, &s.AvailableSeats); err != nil {
			return nil, err
		}
		showtimes = append(showtimes, s)
	}
	return showtimes, rows.Err()
}

func (r *ShowtimeRepo) GetSeats(showtimeID string) ([]Seat, error) {
	rows, err := r.db.Query(
		"SELECT id, row, number, is_available FROM seats WHERE showtime_id=$1 ORDER BY row, number",
		showtimeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var seats []Seat
	for rows.Next() {
		var s Seat
		if err := rows.Scan(&s.ID, &s.Row, &s.Number, &s.IsAvailable); err != nil {
			return nil, err
		}
		s.ShowtimeID = showtimeID
		seats = append(seats, s)
	}
	return seats, rows.Err()
}
