package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

func main() {
	rClient := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
		DB:      0,
	})
	defer rClient.Close()

	locker, ctx := redislock.New(rClient), context.Background()

	lock, err := locker.Obtain(ctx, "event_name", 1*time.Second, nil)
	if errors.Is(err, redislock.ErrNotObtained) {
		fmt.Printf("cannot obtain the locker\n")
	} else if err != nil {
		fmt.Printf("obtain err: %v\n", err)
		return
	}
	defer lock.Release(ctx)

	fmt.Printf("lock success\n")

	time.Sleep(800 * time.Millisecond)
	if ttl, err := lock.TTL(ctx); err != nil {
		fmt.Printf("TTL err: %v\n", err)
		return
	} else if ttl > 0 {
		fmt.Printf("still in lock\n")
	}

	time.Sleep(100 * time.Millisecond)

	if err = lock.Refresh(ctx, 200*time.Millisecond, nil); err != nil {
		fmt.Printf("extend lock err: %v\n", err)
		return
	}
	time.Sleep(200 * time.Millisecond)
	if ttl, err := lock.TTL(ctx); err != nil {
		fmt.Printf("second ttl err: %v\n", err)
		return
	} else if ttl == 0 {
		fmt.Printf("lock expired!\n")
	}

}
