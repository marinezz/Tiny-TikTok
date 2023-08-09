package handler

import (
	"api_router/internal/service"
	"api_router/pkg/auth"
	"api_router/pkg/res"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	var userInfoReq service.UserInfoRequest

	// Todo 两种方法获取参数，没懂为什么get也可以获取   只实现了用户服务相关的其余服务还未实现
	userId := ctx.Query("user_id")
	userInfoReq.UserId, _ = strconv.ParseInt(userId, 10, 64)
	//userId, _ := ctx.Get("user_id")
	//userInfoReq.UserId, _ = userId.(int64)

	userServiceClient := ctx.Keys["user_service"].(service.UserServiceClient)
	userResp, _ := userServiceClient.UserInfo(context.Background(), &userInfoReq)

	r := res.UserInfoResponse{
		StatusCode: userResp.StatusCode,
		StatusMsg:  userResp.StatusMsg,
		User:       userResp.User,
	}

	ctx.JSON(http.StatusOK, r)
}
