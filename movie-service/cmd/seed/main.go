package main

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://movie:movie@localhost:5434/moviedb?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	defer db.Close()

	movies := []struct{ title, genre string; dur int }{
		{"The Shawshank Redemption", "Drama", 142},
		{"The Godfather", "Crime", 175},
		{"The Dark Knight", "Action", 152},
		{"Pulp Fiction", "Crime", 154},
		{"Schindler's List", "History", 195},
		{"Forrest Gump", "Drama", 142},
		{"Inception", "Sci-Fi", 148},
		{"Fight Club", "Drama", 139},
		{"The Matrix", "Sci-Fi", 136},
		{"Interstellar", "Sci-Fi", 169},
		{"Parasite", "Thriller", 132},
		{"Gladiator", "Action", 155},
		{"The Lion King", "Animation", 88},
		{"Avengers: Endgame", "Action", 181},
		{"Joker", "Crime", 122},
		{"Spirited Away", "Animation", 125},
		{"The Silence of the Lambs", "Thriller", 118},
		{"Se7en", "Crime", 127},
		{"The Prestige", "Mystery", 130},
		{"Whiplash", "Drama", 106},
		{"The Wolf of Wall Street", "Comedy", 180},
		{"Dune", "Sci-Fi", 155},
		{"The Truman Show", "Comedy", 103},
		{"Coco", "Animation", 105},
		{"Top Gun: Maverick", "Action", 130},
		{"La La Land", "Musical", 128},
		{"The Green Mile", "Drama", 189},
		{"Goodfellas", "Crime", 146},
		{"Saving Private Ryan", "War", 169},
		{"The Departed", "Crime", 151},
		{"American History X", "Drama", 119},
		{"Memento", "Mystery", 113},
		{"Eternal Sunshine of the Spotless Mind", "Romance", 108},
		{"The Grand Budapest Hotel", "Comedy", 99},
		{"Mad Max: Fury Road", "Action", 120},
		{"No Country for Old Men", "Thriller", 122},
		{"Oldboy", "Mystery", 120},
		{"Snatch", "Comedy", 102},
		{"The Intouchables", "Comedy", 112},
		{"A Beautiful Mind", "Biography", 135},
		{"Life Is Beautiful", "Comedy", 116},
		{"The Usual Suspects", "Crime", 106},
		{"Back to the Future", "Adventure", 116},
		{"Jurassic Park", "Adventure", 127},
		{"The Lord of the Rings: The Return of the King", "Fantasy", 201},
		{"Star Wars: Episode V", "Sci-Fi", 124},
		{"The Social Network", "Biography", 120},
		{"Good Will Hunting", "Drama", 126},
		{"The Shining", "Horror", 146},
		{"Alien", "Horror", 117},
	}
	for _, m := range movies {
		_, err := db.Exec("INSERT INTO movies (title, genre, duration_min, poster_url) VALUES ($1, $2, $3, '')",
			m.title, m.genre, m.dur)
		if err != nil {
			log.Printf("seed error: %v", err)
		}
	}
	log.Println("Seed completed: 50 movies inserted")
}
