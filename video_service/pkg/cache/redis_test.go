package cache

import (
	"context"
	"testing"
)

func TestInitRedis(t *testing.T) {
	InitRedis()
	Redis.Set(context.Background(), "K2", "myValue", 0)
}
