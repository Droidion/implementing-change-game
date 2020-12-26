package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/rotisserie/eris"
	"os"
)

var ctx = context.Background()

var Redis *redis.Client

// InitRedis creates and tests connection to Redis database
func InitRedis() error {
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		return eris.New("could not find non empty REDIS_DSN in environment variables")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return eris.Wrap(err, "could not find non empty REDIS_DSN in environment variables")
	}

	Redis = client
	return nil
}