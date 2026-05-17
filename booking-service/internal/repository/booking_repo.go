package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type Booking struct {
	ID         string
	UserID     string
	ShowtimeID string
	Status     string
}

type BookingRepo struct {
	db *sql.DB
}

func NewBookingRepo(db *sql.DB) *BookingRepo {
	return &BookingRepo{db: db}
}

func (r *BookingRepo) CreateBooking(tx *sql.Tx, userID, showtimeID string, seatIDs []string) (*Booking, error) {
	bookingID := uuid.New().String()
	_, err := tx.Exec(`INSERT INTO bookings (id, user_id, showtime_id, status) VALUES ($1, $2, $3, 'CONFIRMED')`,
		bookingID, userID, showtimeID)
	if err != nil {
		return nil, fmt.Errorf("insert booking: %w", err)
	}
	for _, seatID := range seatIDs {
		_, err := tx.Exec(`INSERT INTO booking_seats (booking_id, seat_id) VALUES ($1, $2)`, bookingID, seatID)
		if err != nil {
			return nil, fmt.Errorf("insert seat: %w", err)
		}
		// блокируем место (меняем is_available = false) – предполагаем, что таблица seats в другой БД, но здесь просто для примера обновим
		_, err = tx.Exec(`UPDATE seats SET is_available = false WHERE id = $1`, seatID)
		if err != nil {
			return nil, fmt.Errorf("update seat: %w", err)
		}
	}
	return &Booking{ID: bookingID, UserID: userID, ShowtimeID: showtimeID, Status: "CONFIRMED"}, nil
}

func (r *BookingRepo) GetBooking(id string) (*Booking, []string, error) {
	row := r.db.QueryRow(`SELECT id, user_id, showtime_id, status FROM bookings WHERE id=$1`, id)
	var b Booking
	err := row.Scan(&b.ID, &b.UserID, &b.ShowtimeID, &b.Status)
	if err != nil {
		return nil, nil, fmt.Errorf("booking not found: %w", err)
	}
	rows, err := r.db.Query(`SELECT seat_id FROM booking_seats WHERE booking_id=$1`, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var seats []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, nil, err
		}
		seats = append(seats, s)
	}
	return &b, seats, rows.Err()
}

func (r *BookingRepo) CancelBooking(tx *sql.Tx, bookingID string) error {
	_, err := tx.Exec(`UPDATE bookings SET status='CANCELLED' WHERE id=$1`, bookingID)
	if err != nil {
		return err
	}
	// возвращаем места
	rows, err := tx.Query(`SELECT seat_id FROM booking_seats WHERE booking_id=$1`, bookingID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var seatID string
		if err := rows.Scan(&seatID); err != nil {
			return err
		}
		_, err = tx.Exec(`UPDATE seats SET is_available = true WHERE id=$1`, seatID)
		if err != nil {
			return err
		}
	}
	return rows.Err()
}
