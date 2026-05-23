package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// InitRedis connects to the Redis instance and verifies readiness with a Ping
func InitRedis() (*redis.Client, error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	password := os.Getenv("REDIS_PASSWORD")

	dbStr := os.Getenv("REDIS_DB")
	dbNum := 0
	if dbStr != "" {
		if parsed, err := strconv.Atoi(dbStr); err == nil {
			dbNum = parsed
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbNum,
	})

	// Enforce a strict 3-second timeout for the initial health check
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Ping Redis to ensure the password and address are correct
	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return rdb, nil
}