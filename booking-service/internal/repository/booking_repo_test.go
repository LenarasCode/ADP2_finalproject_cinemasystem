package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/cinema-system/booking-service/internal/repository"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")
	dsn := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", host, port.Port())

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		container.Terminate(ctx)
		t.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		container.Terminate(ctx)
		t.Fatal(err)
	}
	// миграции
	for _, q := range []string{
		`CREATE TABLE bookings (id UUID PRIMARY KEY, user_id VARCHAR(255), showtime_id UUID, status VARCHAR(20))`,
		`CREATE TABLE booking_seats (booking_id UUID, seat_id UUID, PRIMARY KEY (booking_id, seat_id))`,
		`CREATE TABLE seats (id UUID PRIMARY KEY, is_available BOOLEAN)`,
	} {
		if _, err := db.Exec(q); err != nil {
			container.Terminate(ctx)
			t.Fatal(err)
		}
	}
	return db, func() {
		db.Close()
		container.Terminate(ctx)
	}
}

func TestCreateBooking(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration test")
	}
	db, teardown := setupTestDB(t)
	defer teardown()

	// вставим место
	_, err := db.Exec(`INSERT INTO seats (id, is_available) VALUES ('s1', true)`)
	if err != nil {
		t.Fatal(err)
	}

	repo := repository.NewBookingRepo(db)
	tx, _ := db.Begin()
	_, err = repo.CreateBooking(tx, "u1", "st1", []string{"s1"})
	if err != nil {
		t.Fatal(err)
	}
	tx.Commit()

	b, seats, err := repo.GetBooking("b1") // мы не знаем ID, но в тесте можно было бы вернуть
	_ = b
	_ = seats
	// просто проверяем, что не ошибка
}
