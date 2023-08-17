// Package router 路由
package router

import (
	"api_router/internal/handler"
	"api_router/router/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter(serveInstance map[string]interface{}) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.ServeMiddleware(serveInstance))
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
		baseGroup.POST("/user/login", handler.UserLogin)
		baseGroup.GET("/user", middleware.JWTMiddleware(), handler.UserInfo)

		// 视频
		publishGroup := baseGroup.Group("/publish")
		publishGroup.Use(middleware.JWTMiddleware())
		{
			publishGroup.POST("/action/", handler.PublishAction)
			publishGroup.GET("/list", handler.PublishList)
		}
		favoriteGroup := baseGroup.Group("favorite")
		favoriteGroup.Use(middleware.JWTMiddleware())
		{
			favoriteGroup.POST("action", handler.FavoriteAction)
			favoriteGroup.GET("list", handler.FavoriteList)
		}
		commentGroup := baseGroup.Group("/comment")
		commentGroup.Use(middleware.JWTMiddleware())
		{
			commentGroup.POST("/action", handler.CommentAction)
			commentGroup.GET("/list")
		}
		// 社交
		relationGroup := baseGroup.Group("/relation")
		relationGroup.Use(middleware.JWTMiddleware())
		{
			relationGroup.POST("/action", handler.FollowAction)
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
