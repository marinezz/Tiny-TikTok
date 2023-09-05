package main

import (
	"video/config"
	"video/discovery"
	"video/internal/handler"
	"video/internal/model"
	"video/pkg/cache"
)

func main() {
	config.InitConfig() // 初始话配置文件
	model.InitDb()      // 初始化数据库
	cache.InitRedis()   // 初始化缓
	go func() {
		handler.PublishVideo()
	}()
	discovery.AutoRegister() // 自动注册
}
