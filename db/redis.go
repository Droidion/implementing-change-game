package db

import (
	"context"
	"github.com/Droidion/implementing-change-game/models"
	"github.com/go-redis/redis/v8"
	"github.com/rotisserie/eris"
	"os"
	"strconv"
	"time"
)

var Ctx = context.Background()

// Redis connection
var Redis *redis.Client

// Predicate to app-related Redis keys
const tokenKeyPredicate = "impl-change-token-"

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

	_, err := client.Ping(Ctx).Result()
	if err != nil {
		return eris.Wrap(err, "could not find non empty REDIS_DSN in environment variables")
	}

	Redis = client
	return nil
}

// SaveTokenToRedis saves a single token to Redis
func SaveTokenToRedis(expires int64, uuid string, userId uint64) error {
	utcTime := time.Unix(expires, 0)
	now := time.Now()
	err := Redis.Set(Ctx, tokenKeyPredicate+uuid, strconv.Itoa(int(userId)), utcTime.Sub(now)).Err()
	if err != nil {
		return eris.Wrap(err, "could not save token to Redis")
	}
	return nil
}

// FetchAuth tries to get user id from cached token metadata
func FetchAuth(authD *models.AccessDetails) (uint64, error) {
	userid, err := Redis.Get(Ctx, tokenKeyPredicate+authD.AccessUuid).Result()
	if err != nil {
		return 0, eris.Wrap(err, "could not find token in Redis")
	}

	userID, err := strconv.ParseUint(userid, 10, 64)
	if err != nil {
		return 0, eris.Wrap(err, "could not parse user id")
	}

	return userID, nil
}

// DeleteAuth tries to delete token in redis cache
func DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := Redis.Del(Ctx, tokenKeyPredicate+givenUuid).Result()
	if err != nil {
		return 0, eris.Wrap(err, "could not find token to delete in Redis cache")
	}
	return deleted, nil
}
