package handler

import (
	"api_router/internal/service"
	"api_router/pkg/auth"
	"api_router/pkg/res"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserRegister 用户注册
func UserRegister(ctx *gin.Context) {
	var userReq service.UserRequest
	ctx.Bind(&userReq)

	userServiceClient := ctx.Keys["user_service"].(service.UserServiceClient)
	userResp, err := userServiceClient.UserRegister(context.Background(), &userReq)
	if err != nil {
		panic(err)
	}
	token, _ := auth.GenerateToken(userResp.UserId)

	r := res.UserResponse{
		StatusCode: userResp.StatusCode,
		StatusMsg:  userResp.StatusMsg,
		UserId:     userResp.UserId,
		Token:      token,
	}

	ctx.JSON(http.StatusOK, r)
}

func UserLogin(ctx *gin.Context) {
	var userReq service.UserRequest
	ctx.Bind(&userReq)

	userServiceClient := ctx.Keys["user_service"].(service.UserServiceClient)
	userResp, err := userServiceClient.UserLogin(context.Background(), &userReq)
	if err != nil {
		panic(err)
	}
	token, _ := auth.GenerateToken(userResp.UserId)

	r := res.UserResponse{
		StatusCode: userResp.StatusCode,
		StatusMsg:  userResp.StatusMsg,
		UserId:     userResp.UserId,
		Token:      token,
	}

	ctx.JSON(http.StatusOK, r)
}

func UserInfo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"code":    1011,
	})
}
