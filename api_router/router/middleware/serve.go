// 服务中间件，接收服务实例，并保存到context.Key中

package middleware

import "github.com/gin-gonic/gin"

func ServeMiddleware(service []interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Todo 目前想到的就是顺序存储，且不能打乱顺序，看以后有更好的方法没

		c.Next()
	}
}
