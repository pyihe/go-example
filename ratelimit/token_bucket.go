package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

// 令牌桶算法：算法以固定速率往桶里生成令牌，直到桶满。
// 如果请求频繁，超过设定的频率，当桶里的令牌被取完后，令牌桶算法将会转化为漏桶算法，即以固定频率执行API
// 如果API执行速率低于设定值，则在令牌被取完前，令牌桶都可以抵挡某一瞬间迸发的一定量的高频请求

func main() {

	// 表示100/time.Second， 桶里最多装50个
	var limiter = rate.NewLimiter(10, 50)
	var before = time.Now()
	for {
		if !limiter.Allow() {
			if err := limiter.Wait(context.Background()); err != nil {
				fmt.Println(err)
				break
			}
		}
		now := time.Now()
		gap := now.Sub(before)
		before = now
		fmt.Printf("%s\n", gap.String())
	}
}
