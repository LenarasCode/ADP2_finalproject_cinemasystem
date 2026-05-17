package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/cinema-system/showtime-service/internal/repository"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*repository.ShowtimeRepo, func()) {
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

	for _, q := range []string{
		`CREATE TABLE showtimes (id UUID PRIMARY KEY, movie_id VARCHAR(255) NOT NULL, hall VARCHAR(100) NOT NULL, start_time TIMESTAMPTZ NOT NULL, available_seats INT NOT NULL)`,
		`CREATE TABLE seats (id UUID PRIMARY KEY, showtime_id UUID NOT NULL REFERENCES showtimes(id), row VARCHAR(5) NOT NULL, number INT NOT NULL, is_available BOOLEAN DEFAULT true)`,
	} {
		if _, err := db.Exec(q); err != nil {
			container.Terminate(ctx)
			t.Fatal(err)
		}
	}

	repo := repository.NewShowtimeRepo(db)
	return repo, func() {
		db.Close()
		container.Terminate(ctx)
	}
}

func TestGetShowtime_Integration(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration test; set INTEGRATION=1 to run")
	}
	repo, teardown := setupTestDB(t)
	defer teardown()

	testID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	_, err := repo.DB.Exec(`INSERT INTO showtimes (id, movie_id, hall, start_time, available_seats) VALUES ($1, 'm1', 'H1', '2025-01-01T10:00:00Z', 80)`, testID)
	if err != nil {
		t.Fatal(err)
	}

	s, err := repo.GetShowtime(testID)
	if err != nil {
		t.Fatal(err)
	}
	if s.ID != testID || s.AvailableSeats != 80 {
		t.Errorf("unexpected showtime: %+v", s)
	}
}
