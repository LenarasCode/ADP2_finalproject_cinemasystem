package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/cinema-system/showtime-service/internal/repository"
)

type ShowtimeCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewShowtimeCache(client *redis.Client) *ShowtimeCache {
	return &ShowtimeCache{client: client, ttl: 5 * time.Minute}
}

func (c *ShowtimeCache) GetShowtime(ctx context.Context, id string) (*repository.Showtime, error) {
	val, err := c.client.Get(ctx, "showtime:"+id).Bytes()
	if err != nil {
		return nil, err
	}
	var s repository.Showtime
	if err := json.Unmarshal(val, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *ShowtimeCache) SetShowtime(ctx context.Context, s *repository.Showtime) error {
	data, _ := json.Marshal(s)
	return c.client.Set(ctx, "showtime:"+s.ID, data, c.ttl).Err()
}

func (c *ShowtimeCache) DeleteShowtime(ctx context.Context, id string) error {
	return c.client.Del(ctx, "showtime:"+id).Err()
}
