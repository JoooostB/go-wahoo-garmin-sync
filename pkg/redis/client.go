package redis

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func Default() (*redis.Client, context.Context, error) {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "", // no password set
		DB:       0,  // use default DB
	}), ctx, nil
}

// Set key (k) & value (v)
func Set(rdb *redis.Client, k, v string) error {
	err := rdb.Set(ctx, k, v, 0).Err()
	if err != nil {
		panic(err)
	}
	return nil
}
