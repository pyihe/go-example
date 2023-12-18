package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	_ = rdb.FlushDB(ctx).Err()

	limiter := redis_rate.NewLimiter(rdb)
	res, err := limiter.AllowN(ctx, "key1", redis_rate.PerSecond(10), 2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("allowed: %v, remaining: %v\n", res.Allowed, res.Remaining)
}
