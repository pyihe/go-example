package main

import (
	"fmt"
	"time"

	"go.uber.org/ratelimit"
)

// 漏桶算法：算法以固定速率限定API的执行，在请求频繁的情况下，API最多只能按照令牌生成的速率被执行，比如100/time.Second表示每秒执行100次，
// 每次间隔10ms左右

func main() {
	// 表示100/time.Second
	var limiter = ratelimit.New(100)
	var before = time.Now()

	for {
		limiter.Take()
		now := time.Now()
		gap := now.Sub(before)
		before = now
		fmt.Printf("%s\n", gap.String())
	}
}
