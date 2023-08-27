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

// FavoriteAction 喜欢操作
func FavoriteAction(ctx *gin.Context) {
	var favoriteActionReq service.FavoriteActionRequest

	userId, _ := ctx.Get("user_id")
	favoriteActionReq.UserId, _ = userId.(int64)
	// string转int64
	videoId := ctx.PostForm("video_id")
	if videoId == "" {
		videoId = ctx.Query("video_id")
	}
	favoriteActionReq.VideoId, _ = strconv.ParseInt(videoId, 10, 64)

	actionType := ctx.PostForm("action_type")
	if actionType == "" {
		actionType = ctx.Query("action_type")
	}
	actionTypeValue, _ := strconv.Atoi(actionType)

	// 异常操作
	if actionTypeValue == 1 || actionTypeValue == 2 {
		favoriteActionReq.ActionType = int64(actionTypeValue)

		videoServiceClient := ctx.Keys["video_service"].(service.VideoServiceClient)
		videoServiceResp, err := videoServiceClient.FavoriteAction(context.Background(), &favoriteActionReq)
		if err != nil {
			PanicIfFavoriteError(err)
		}

		r := res.FavoriteActionResponse{
			StatusCode: videoServiceResp.StatusCode,
			StatusMsg:  videoServiceResp.StatusMsg,
		}

		ctx.JSON(http.StatusOK, r)
	} else {
		r := res.FavoriteActionResponse{
			StatusCode: exception.ErrOperate,
			StatusMsg:  exception.GetMsg(exception.ErrOperate),
		}

		ctx.JSON(http.StatusOK, r)
	}
}

func FavoriteList(ctx *gin.Context) {
	var favoriteListReq service.FavoriteListRequest

	userIdStr := ctx.Query("user_id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	favoriteListReq.UserId = userId

	videoServiceClient := ctx.Keys["video_service"].(service.VideoServiceClient)
	favoriteListResp, err := videoServiceClient.FavoriteList(context.Background(), &favoriteListReq)
	if err != nil {
		PanicIfFavoriteError(err)
	}

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
