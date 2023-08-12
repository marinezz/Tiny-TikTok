package main

import (
	"video/config"
	"video/discovery"
	"video/internal/model"
)

func main() {
	config.InitConfig()      // 初始话配置文件
	model.InitDb()           // 初始化数据库
	discovery.AutoRegister() // 自动注册
}
