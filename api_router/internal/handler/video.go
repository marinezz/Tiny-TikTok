package handler

import (
	"api_router/internal/service"
	"api_router/pkg/res"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Feed(ctx *gin.Context) {
	var feedReq service.FeedRequest

	token, _ := ctx.Get("token")
	if token == "" {
		feedReq.UserId = -1
	} else {
		userId, _ := ctx.Get("user_id")
		feedReq.UserId, _ = userId.(int64)
	}

	latestTime := ctx.Query("latest_time")
	timePoint, _ := time.Parse("2006-01-02 15:04:05", latestTime)
	feedReq.LatestTime = timePoint.Unix()

	videoServiceClient := ctx.Keys["video_service"].(service.VideoServiceClient)
	feedResp, _ := videoServiceClient.Feed(context.Background(), &feedReq)

	var userIds []int64
	for _, video := range feedResp.VideoList {
		userIds = append(userIds, video.AuthId)
	}

	// 找到所有的用户信息
	userInfos := GetUserInfo(userIds, ctx)

	list := BuildVideoList(feedResp.VideoList, userInfos)

	r := res.FeedResponse{
		StatusCode: feedResp.StatusCode,
		StatusMsg:  feedResp.StatusMsg,
		NextTime:   feedResp.NextTime,
		VideoList:  list,
	}

	ctx.JSON(http.StatusOK, r)
}

// PublishAction 发布视频
func PublishAction(ctx *gin.Context) {
	var publishActionReq service.PublishActionRequest

	userId, _ := ctx.Get("user_id")
	publishActionReq.UserId = userId.(int64)

	publishActionReq.Title = ctx.PostForm("title")

	formFile, _ := ctx.FormFile("data")
	file, err := formFile.Open()
	if err != nil {
		PanicIfPublishError(err)
	}
	defer file.Close()
	buf, err := io.ReadAll(file) // 将文件读取到字节切片buf中
	if err != nil {
		PanicIfPublishError(err)
	}
	publishActionReq.Data = buf

	videoServiceClient := ctx.Keys["video_service"].(service.VideoServiceClient)
	videoServiceResp, _ := videoServiceClient.PublishAction(context.Background(), &publishActionReq)

	r := res.PublishActionResponse{
		StatusCode: videoServiceResp.StatusCode,
		StatusMsg:  videoServiceResp.StatusMsg,
	}

	ctx.JSON(http.StatusOK, r)
}

// PublishList 发布列表
func PublishList(ctx *gin.Context) {
	var pulishListReq service.PublishListRequest

	userIdStr := ctx.Query("user_id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	pulishListReq.UserId = userId

	videoServiceClient := ctx.Keys["video_service"].(service.VideoServiceClient)
	publishListResp, _ := videoServiceClient.PublishList(context.Background(), &pulishListReq)

	var userIds []int64
	for _, video := range publishListResp.VideoList {
		userIds = append(userIds, video.AuthId)
	}

	// 找到所有的用户信息
	userInfos := GetUserInfo(userIds, ctx)

	list := BuildVideoList(publishListResp.VideoList, userInfos)

	r := res.VideoListResponse{
		StatusCode: publishListResp.StatusCode,
		StatusMsg:  publishListResp.StatusMsg,
		VideoList:  list,
	}

	ctx.JSON(http.StatusOK, r)
}

// BuildVideoList 构建视频列表
func BuildVideoList(videos []*service.Video, userInfos []res.User) []res.Video {

	var videoList []res.Video

	for i, video := range videos {
		videoList = append(videoList, res.Video{
			Id:            video.Id,
			Author:        userInfos[i],
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavorite,
			Title:         video.Title,
		})
	}

	return videoList
}

// PanicIfPublishError 错误处理
func PanicIfPublishError(err error) {
	if err != nil {
		err = errors.New("publishService--error--" + err.Error())
		// Todo 统一的日志处理
		log.Print(err)
		panic(err)
	}
}
