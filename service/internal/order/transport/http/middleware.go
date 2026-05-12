package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimiterMiddleware(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		key := fmt.Sprintf("rate_limit:%s", c.ClientIP())

		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if count == 1 {
			rdb.Expire(ctx, key, window)
		}

		if int(count) > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			return
		}

		c.Next()
	}
}
