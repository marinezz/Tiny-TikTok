// Package middleware  token校验中间件
package middleware

import (
	"api_router/pkg/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"utils/exception"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = 200
		// token 可能在query中也可能在postForm中
		token := c.Query("token")
		if token == "" {
			token = c.PostForm("token")
		}

		// token不存在
		if token == "" {
			code = exception.RequestERROR
		}

		// 验证token（验证不通过，或者超时）
		claims, err := auth.ParseToken(token)
		if err != nil {
			code = exception.UnAuth
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = exception.TokenTimeOut
		}

		if code != exception.SUCCESS {
			c.JSON(http.StatusOK, gin.H{
				"StatusCode": code,
				"StatusMsg":  exception.GetMsg(code),
			})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserId)
		c.Next()
	}
}
