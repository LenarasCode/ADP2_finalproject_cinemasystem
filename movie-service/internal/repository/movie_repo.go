package repository

import (
	"database/sql"
	"fmt"
)

type Movie struct {
	ID          string
	Title       string
	Genre       string
	DurationMin int
	PosterURL   string
}

type MovieRepo struct {
	db *sql.DB
}

func NewMovieRepo(db *sql.DB) *MovieRepo {
	return &MovieRepo{db: db}
}

func (r *MovieRepo) GetMovie(id string) (*Movie, error) {
	row := r.db.QueryRow("SELECT id, title, genre, duration_min, poster_url FROM movies WHERE id=$1", id)
	var m Movie
	err := row.Scan(&m.ID, &m.Title, &m.Genre, &m.DurationMin, &m.PosterURL)
	if err != nil {
		return nil, fmt.Errorf("movie not found: %w", err)
	}
	return &m, nil
}

func (r *MovieRepo) ListMovies(page, limit int) ([]Movie, int, error) {
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM movies").Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	rows, err := r.db.Query(
		"SELECT id, title, genre, duration_min, poster_url FROM movies ORDER BY title LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var movies []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Genre, &m.DurationMin, &m.PosterURL); err != nil {
			return nil, 0, err
		}
		movies = append(movies, m)
	}
	return movies, total, rows.Err()
}

func (r *MovieRepo) SearchMovies(query string) ([]Movie, error) {
	rows, err := r.db.Query(
		"SELECT id, title, genre, duration_min, poster_url FROM movies WHERE title ILIKE $1 OR genre ILIKE $1",
		"%"+query+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var movies []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Genre, &m.DurationMin, &m.PosterURL); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, rows.Err()
}
