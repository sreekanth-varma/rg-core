package rgcache

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/sreekanth-varma/rg-core/rgutil"
)

var (
	rdb *redis.Client
)

func Init() rgutil.Err {

	// Get Redis host and password from environment variables
	host := os.Getenv("redis_url")
	if host == "" {
		log.Fatalln("cache: redis_host cannot be empty")
		return rgutil.ErrProcessingFailed
	}

	password := os.Getenv("redis_password")

	// Initialize Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Println("cache: connected succcessfully")
			return nil
		},
	})

	if ok := CheckHealth(); !ok {
		slog.Error("rgcache: connection failed")
		panic("cache: connection failed")
	}

	return rgutil.ErrNil
}

func CheckHealth() bool {
	_, err := rdb.ClientID(context.Background()).Result()
	if err != nil {
		slog.Error("rgcache: connection failed", "error", err)
		return false
	}

	return true
}
