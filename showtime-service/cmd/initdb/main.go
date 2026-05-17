package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("SUPERUSER_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	defer db.Close()

	// Создаём базу
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'showtimedb')").Scan(&exists)
	if err != nil {
		log.Fatalf("check db: %v", err)
	}
	if !exists {
		if _, err := db.Exec("CREATE DATABASE showtimedb"); err != nil {
			log.Fatalf("create db: %v", err)
		}
		log.Println("Database showtimedb created")
	} else {
		log.Println("Database showtimedb already exists")
	}

	// Создаём пользователя и даём права
	_, err = db.Exec(`
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'showtime') THEN
			CREATE USER showtime WITH PASSWORD 'showtime';
		END IF;
	END
	$$;
	GRANT ALL PRIVILEGES ON DATABASE showtimedb TO showtime;
	`)
	if err != nil {
		log.Fatalf("create user: %v", err)
	}
	log.Println("User showtime ready")
}
