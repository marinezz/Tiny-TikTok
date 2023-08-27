package handler

import (
	"api_router/internal/service"
	"api_router/pkg/res"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func PostMessage(ctx *gin.Context) {
	var postMessage service.PostMessageRequest
	userId, _ := ctx.Get("user_id")
	postMessage.UserId, _ = userId.(int64)
	toUserId := ctx.PostForm("to_user_id")
	if toUserId == "" {
		toUserId = ctx.Query("to_user_id")
	}
	postMessage.ToUserId, _ = strconv.ParseInt(toUserId, 10, 64)
	actionType := ctx.PostForm("action_type")
	if actionType == "" {
		actionType = ctx.Query("action_type")
	}
	actionTypeInt64, _ := strconv.ParseInt(actionType, 10, 32)
	postMessage.ActionType = int32(actionTypeInt64)
	content := ctx.PostForm("content")
	if content == "" {
		content = ctx.Query("content")
	}
	postMessage.Content = content

	socialServiceClient := ctx.Keys["social_service"].(service.SocialServiceClient)
	socialResp, err := socialServiceClient.PostMessage(context.Background(), &postMessage)
	if err != nil {
		PanicIfMessageError(err)
	}

	r := res.PostMessageResponse{
		StatusCode: socialResp.StatusCode,
		StatusMsg:  socialResp.StatusMsg,
	}
	ctx.JSON(http.StatusOK, r)
}

func GetMessage(ctx *gin.Context) {
	var getMessage service.GetMessageRequest
	userId, _ := ctx.Get("user_id")
	getMessage.UserId, _ = userId.(int64)
	toUserId := ctx.Query("to_user_id")
	getMessage.ToUserId, _ = strconv.ParseInt(toUserId, 10, 64)

	socialServiceClient := ctx.Keys["social_service"].(service.SocialServiceClient)
	socialResp, err := socialServiceClient.GetMessage(context.Background(), &getMessage)

	if err != nil {
		PanicIfMessageError(err)
	}

	r := new(res.GetMessageResponse)
	r.StatusCode = socialResp.StatusCode
	r.StatusMsg = socialResp.StatusMsg
	for _, message := range socialResp.Message {
		createTime, _ := time.Parse("2006-01-02 15:04:05", message.CreatedAt)
		messageResp := res.Message{
			Id:         message.Id,
			ToUserId:   message.ToUserId,
			FromUserID: message.UserId,
			Content:    message.Content,
			CreateTime: createTime.Unix(),
		}
		r.MessageList = append(r.MessageList, messageResp)
	}

	ctx.JSON(http.StatusOK, r)
}
