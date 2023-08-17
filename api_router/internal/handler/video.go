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
)

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
	token := ctx.Query("token")
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
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
