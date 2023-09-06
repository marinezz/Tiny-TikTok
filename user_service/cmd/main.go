package main

import (
	"user/config"
	"user/discovery"
	"user/internal/model"
	"user/pkg/cache"
)

func main() {
	config.InitConfig()      // 初始话配置文件
	cache.InitRedis()        // 初始化redis
	model.InitDb()           // 初始化数据库
	discovery.AutoRegister() // 自动注册
}
