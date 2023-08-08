package main

import (
	"user/config"
	"user/discovery"
	"user/internal/model"
)

func main() {
	config.InitConfig()      // 初始话配置文件
	model.InitDb()           // 初始化数据库
	discovery.AutoRegister() // 自动注册
}
