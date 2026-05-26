package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/silasrm/api-obm/internal/infrastructure/config"
)

type CacheRepo struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCacheRepo(cfg config.RedisConfig) *CacheRepo {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v. Cache will be disabled.", err)
		return &CacheRepo{client: nil, ttl: time.Duration(cfg.CacheTTL) * time.Hour}
	}

	return &CacheRepo{
		client: client,
		ttl:    time.Duration(cfg.CacheTTL) * time.Hour,
	}
}

func (r *CacheRepo) Get(ctx context.Context, key string) (string, error) {
	if r.client == nil {
		return "", nil
	}

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *CacheRepo) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if r.client == nil {
		return nil
	}

	if ttl <= 0 {
		ttl = r.ttl
	}

	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *CacheRepo) DeleteByPattern(ctx context.Context, pattern string) error {
	if r.client == nil {
		return nil
	}

	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (r *CacheRepo) HealthCheck(ctx context.Context) error {
	if r.client == nil {
		return fmt.Errorf("redis not available")
	}

	return r.client.Ping(ctx).Err()
}

func (r *CacheRepo) Close() error {
	if r.client == nil {
		return nil
	}

	return r.client.Close()
}
