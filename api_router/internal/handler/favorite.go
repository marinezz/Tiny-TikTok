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
	favoriteActionReq.ActionType = int32(actionTypeValue)

	videoServiceClient := ctx.Keys["video_service"].(service.VideoServiceClient)
	videoServiceResp, _ := videoServiceClient.FavoriteAction(context.Background(), &favoriteActionReq)

	r := res.FavoriteActionResponse{
		StatusCode: videoServiceResp.StatusCode,
		StatusMsg:  videoServiceResp.StatusMsg,
	}

	ctx.JSON(http.StatusOK, r)
}
