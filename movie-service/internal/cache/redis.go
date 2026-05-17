package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/cinema-system/movie-service/internal/repository"
)

type MovieCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewMovieCache(client *redis.Client) *MovieCache {
	return &MovieCache{client: client, ttl: 10 * time.Minute}
}

func (c *MovieCache) GetMovie(ctx context.Context, id string) (*repository.Movie, error) {
	val, err := c.client.Get(ctx, "movie:"+id).Bytes()
	if err != nil {
		return nil, err
	}
	var m repository.Movie
	if err := json.Unmarshal(val, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (c *MovieCache) SetMovie(ctx context.Context, m *repository.Movie) error {
	data, _ := json.Marshal(m)
	return c.client.Set(ctx, "movie:"+m.ID, data, c.ttl).Err()
}
