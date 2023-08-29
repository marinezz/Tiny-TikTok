package handler

import (
	"api_router/internal/service"
	"api_router/pkg/res"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"utils/exception"
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
	toUserId := ctx.Query("to_user_id")
	followAction.ToUserId, _ = strconv.ParseInt(toUserId, 10, 64)
	actionType := ctx.Query("action_type")
	actionTypeInt64, _ := strconv.ParseInt(actionType, 10, 32)
	followAction.ActionType = int32(actionTypeInt64)

	if actionTypeInt64 != 1 && actionTypeInt64 != 2 {
		r := res.FavoriteActionResponse{
			StatusCode: exception.ErrOperate,
			StatusMsg:  exception.GetMsg(exception.ErrOperate),
		}

		ctx.JSON(http.StatusOK, r)
		return
	}

	socialServiceClient := ctx.Keys["social_service"].(service.SocialServiceClient)
	socialResp, err := socialServiceClient.FollowAction(context.Background(), &followAction)
	if err != nil {
		PanicIfFollowError(err)
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
		PanicIfFollowError(err)
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
		PanicIfFollowError(err)
	}

	r := res.FollowListResponse{
		StatusCode: socialResp.StatusCode,
		StatusMsg:  socialResp.StatusMsg,
		UserList:   GetUserInfo(socialResp.UserId, ctx),
	}
	ctx.JSON(http.StatusOK, r)
}

func GetFriendList(ctx *gin.Context) {
	var friendList service.FollowListRequest
	userId := ctx.Query("user_id")
	friendList.UserId, _ = strconv.ParseInt(userId, 10, 64)

	socialServiceClient := ctx.Keys["social_service"].(service.SocialServiceClient)
	socialResp, err := socialServiceClient.GetFriendList(context.Background(), &friendList)
	if err != nil {
		PanicIfFollowError(err)
	}

	r := res.FollowListResponse{
		StatusCode: socialResp.StatusCode,
		StatusMsg:  socialResp.StatusMsg,
		UserList:   GetUserInfo(socialResp.UserId, ctx),
	}
	ctx.JSON(http.StatusOK, r)
}
