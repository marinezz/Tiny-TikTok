package handler

import (
	"api_router/internal/service"
	"api_router/pkg/res"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Follow struct {
	IsFollow      bool
	FollowCount   int64
	FollowerCount int64
}

func FollowAction(ctx *gin.Context) {
	var followAction service.FollowRequest
	userId, _ := ctx.Get("user_id")
	followAction.UserId, _ = userId.(int64)
	toUserId := ctx.PostForm("to_user_id")
	followAction.ToUserId, _ = strconv.ParseInt(toUserId, 10, 64)
	actionType := ctx.PostForm("action_type")
	actionTypeInt64, _ := strconv.ParseInt(actionType, 10, 32)
	followAction.ActionType = int32(actionTypeInt64)

	socialServiceClient := ctx.Keys["social_service"].(service.SocialServiceClient)
	socialResp, err := socialServiceClient.FollowAction(context.Background(), &followAction)
	if err != nil {
		panic(err)
	}

	r := res.FollowActionResponse{
		StatusCode: socialResp.StatusCode,
		StatusMsg:  socialResp.StatusMsg,
	}

	ctx.JSON(http.StatusOK, r)
}

func GetFollowList(ctx *gin.Context) {
	var followList service.FollowListRequest
	userId := ctx.Query("user_id")
	followList.UserId, _ = strconv.ParseInt(userId, 10, 64)

	socialServiceClient := ctx.Keys["social_service"].(service.SocialServiceClient)
	socialResp, err := socialServiceClient.GetFollowList(context.Background(), &followList)
	if err != nil {
		panic(err)
	}

	r := res.FollowListResponse{
		StatusCode: socialResp.StatusCode,
		StatusMsg:  socialResp.StatusMsg,
		UserList:   GetUserInfo(socialResp.UserId, ctx),
	}
	ctx.JSON(http.StatusOK, r)
}

func GetFollowerList(ctx *gin.Context) {
	var followerList service.FollowListRequest
	userId := ctx.Query("user_id")
	followerList.UserId, _ = strconv.ParseInt(userId, 10, 64)

	socialServiceClient := ctx.Keys["social_service"].(service.SocialServiceClient)
	socialResp, err := socialServiceClient.GetFollowerList(context.Background(), &followerList)
	if err != nil {
		panic(err)
	}

	r := res.FollowListResponse{
		StatusCode: socialResp.StatusCode,
		StatusMsg:  socialResp.StatusMsg,
		UserList:   GetUserInfo(socialResp.UserId, ctx),
	}
	ctx.JSON(http.StatusOK, r)
}
