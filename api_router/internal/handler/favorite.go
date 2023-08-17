package handler

import (
	"api_router/internal/service"
	"api_router/pkg/res"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FavoriteAction 喜欢操作
func FavoriteAction(ctx *gin.Context) {
	var favoriteActionReq service.FavoriteActionRequest

	userId, _ := ctx.Get("user_id")
	favoriteActionReq.UserId, _ = userId.(int64)
	// string转int64
	videoId := ctx.PostForm("video_id")
	favoriteActionReq.VideoId, _ = strconv.ParseInt(videoId, 10, 64)
	// string转int32
	actionType := ctx.PostForm("action_type")
	actionTypeValue, _ := strconv.Atoi(actionType)
	favoriteActionReq.ActionType = int64(actionTypeValue)

	videoServiceClient := ctx.Keys["video_service"].(service.VideoServiceClient)
	videoServiceResp, _ := videoServiceClient.FavoriteAction(context.Background(), &favoriteActionReq)

	r := res.FavoriteActionResponse{
		StatusCode: videoServiceResp.StatusCode,
		StatusMsg:  videoServiceResp.StatusMsg,
	}

	ctx.JSON(http.StatusOK, r)
}

func FavoriteList(ctx *gin.Context) {
	var favoriteListReq service.FavoriteListRequest

	userIdStr := ctx.Query("user_id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	favoriteListReq.UserId = userId

	videoServiceClient := ctx.Keys["video_service"].(service.VideoServiceClient)
	favoriteListResp, _ := videoServiceClient.FavoriteList(context.Background(), &favoriteListReq)

	// 找到所有的用户Id
	var userIds []int64
	for _, video := range favoriteListResp.VideoList {
		userIds = append(userIds, video.AuthId)
	}

	// 找到所有的用户信息
	userInfos := GetUserInfo(userIds, ctx)

	list := BuildVideoList(favoriteListResp.VideoList, userInfos)

	r := res.VideoListResponse{
		StatusCode: favoriteListResp.StatusCode,
		StatusMsg:  favoriteListResp.StatusMsg,
		VideoList:  list,
	}

	ctx.JSON(http.StatusOK, r)
}
