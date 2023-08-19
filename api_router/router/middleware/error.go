package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"utils/exception"
)

// ErrorMiddleWare 错误处理中间件，捕获panic抛出异常
func ErrorMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			r := recover()
			if r != nil {
				c.JSON(http.StatusOK, gin.H{
					"status_code": exception.ERROR,
					// 打印具体错误
					"status_msg": fmt.Sprintf("%s", r),
				})
				// 中断
				c.Abort()
			}
		}()
		c.Next()
	}
}
