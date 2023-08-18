package main

import (
	"social/config"
	"social/discovery"
	"social/internal/model"
)

func main() {
	config.InitConfig()
	model.InitDb()
	discovery.AutoRegister()
}
