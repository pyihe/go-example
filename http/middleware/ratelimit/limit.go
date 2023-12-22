package ratelimit

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// WithLimiter 限流中间件
// limiter golang.org/x/time/rate 下的令牌桶限流器
// timeout 获取令牌的超时时间，超时即视为获取失败，本次请求不做处理
func WithLimiter(limiter *rate.Limiter, timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, timeout)
		defer cancel()

		if err := limiter.Wait(ctx); err != nil {
			c.AbortWithStatus(http.StatusRequestTimeout)
			return
		}
		c.Next()
	}
}
