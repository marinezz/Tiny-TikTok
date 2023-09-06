package main

import (
	"social/config"
	"social/discovery"
	"social/internal/model"
	"social/pkg/cache"
)

func main() {
	config.InitConfig()
	model.InitDb()
	cache.InitRedis()
	go cache.TimerSync()
	discovery.AutoRegister()
}
