package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://showtime:showtime@localhost:5433/showtimedb?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer db.Close()

	movieIDs := make([]string, 50)
	for i := 0; i < 50; i++ {
		movieIDs[i] = uuid.New().String()
	}

	halls := []string{"Hall 1", "Hall 2", "Hall 3", "Hall 4"}
	now := time.Now()

	for _, movieID := range movieIDs {
		for j := 0; j < 4; j++ {
			showtimeID := uuid.New().String()
			startTime := now.Add(time.Duration(24*j) * time.Hour).Add(time.Duration(10+j) * time.Hour)
			hall := halls[j%len(halls)]

			_, err := db.Exec(
				"INSERT INTO showtimes (id, movie_id, hall, start_time, available_seats) VALUES ($1, $2, $3, $4, 100)",
				showtimeID, movieID, hall, startTime,
			)
			if err != nil {
				log.Fatalf("insert showtime: %v", err)
			}

			for row := 'A'; row <= 'J'; row++ {
				for num := 1; num <= 10; num++ {
					seatID := uuid.New().String()
					_, err := db.Exec(
						"INSERT INTO seats (id, showtime_id, row, number, is_available) VALUES ($1, $2, $3, $4, true)",
						seatID, showtimeID, string(row), num,
					)
					if err != nil {
						log.Fatalf("insert seat: %v", err)
					}
				}
			}
		}
	}
	fmt.Println("Seed completed: 200 showtimes with 100 seats each")
}
