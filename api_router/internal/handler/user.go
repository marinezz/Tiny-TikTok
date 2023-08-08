package handler

import (
	"api_router/internal/service"
	"github.com/gin-gonic/gin"
)

// UserRegister 用户注册
func UserRegister(ctx *gin.Context) {
	var userReq service.UserRequest
	ctx.Bind(&userReq)
}
