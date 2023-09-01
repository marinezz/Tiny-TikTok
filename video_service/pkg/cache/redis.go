package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var Redis *redis.Client

var Ctx = context.Background()

// InitRedis 连接redis
func InitRedis() {
	addr := viper.GetString("redis.address")
	Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0, // 存入DB0
	})
}
