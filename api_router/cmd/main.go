package main

import (
	config "api_router/configs"
	"api_router/discovery"
	"api_router/router"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func main() {
	config.InitConfig()
	resolver := discovery.Resolver()
	r := router.InitRouter(resolver)
	server := &http.Server{
		Addr:           viper.GetString("server.port"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()
}
