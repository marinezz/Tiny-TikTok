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

func PublishList(ctx *gin.Context) {
	token := ctx.Query("token")
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func PanicIfPublishError(err error) {
	if err != nil {
		err = errors.New("publishService--error--" + err.Error())
		// Todo 统一的日志处理
		log.Print(err)
		panic(err)
	}
}
