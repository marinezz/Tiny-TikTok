// Package router 路由
package router

import (
	"api_router/internal/handler"
	"api_router/router/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init() *gin.Engine {
	r := gin.Default()

	// 测试
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	baseGroup := r.Group("/douyin")
	{
		// 视频流
		baseGroup.GET("/feed")

		// 用户
		baseGroup.POST("/user/register", handler.UserRegister)
		baseGroup.POST("/user/login")
		baseGroup.GET("/user", middleware.JWTMiddleware())

		// 视频
		publishGroup := baseGroup.Group("/publish")
		publishGroup.Use(middleware.JWTMiddleware())
		{
			publishGroup.POST("/action")
			publishGroup.GET("/list")
		}
		favoriteGroup := baseGroup.Group("favorite")
		favoriteGroup.Use(middleware.JWTMiddleware())
		{
			favoriteGroup.POST("action")
			favoriteGroup.GET("list")
		}
		commentGroup := baseGroup.Group("/comment")
		commentGroup.Use(middleware.JWTMiddleware())
		{
			commentGroup.POST("/action")
			commentGroup.GET("/list")
		}
		// 社交
		relationGroup := baseGroup.Group("/relation")
		relationGroup.Use(middleware.JWTMiddleware())
		{
			relationGroup.POST("/action")
			relationGroup.GET("/follow/list")
			relationGroup.GET("/follower/list")
			relationGroup.GET("/friend/list")
		}
		messageGroup := baseGroup.Group("/message")
		messageGroup.Use(middleware.JWTMiddleware())
		{
			messageGroup.POST("/action")
			messageGroup.GET("/chat")
		}
	}
	return r
}
