package repository

import (
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis(addr string) {
	RDB = redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
