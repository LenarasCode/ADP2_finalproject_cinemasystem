package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/cinema-system/notification-service/internal/repository"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*repository.EmailLogRepo, func()) {
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
	_, err = db.Exec(`CREATE TABLE email_logs (
		id UUID PRIMARY KEY,
		recipient VARCHAR(255) NOT NULL,
		subject TEXT NOT NULL,
		body TEXT NOT NULL,
		sent_at TIMESTAMPTZ NOT NULL,
		status VARCHAR(50) NOT NULL
	)`)
	if err != nil {
		container.Terminate(ctx)
		t.Fatal(err)
	}
	repo := repository.NewEmailLogRepo(db)
	return repo, func() {
		db.Close()
		container.Terminate(ctx)
	}
}

func TestInsertAndGetAll(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration test")
	}
	repo, teardown := setupTestDB(t)
	defer teardown()

	log := &repository.EmailLog{
		ID: "abc123", Recipient: "x@y.com", Subject: "Test", Body: "Hello", Status: "OK",
	}
	err := repo.Insert(log)
	if err != nil {
		t.Fatal(err)
	}
	logs, err := repo.GetAll()
	if err != nil || len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d, err %v", len(logs), err)
	}
}
