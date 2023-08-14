package handler

import (
	"api_router/internal/service"
	"api_router/pkg/res"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CommentAction 评论操作
func CommentAction(ctx *gin.Context) {
	var commentActionReq service.CommentActionRequest

	userId, _ := ctx.Get("user_id")
	commentActionReq.UserId, _ = userId.(int64)

	videoId := ctx.PostForm("video_id")
	commentActionReq.VideoId, _ = strconv.ParseInt(videoId, 10, 64)

	actionType := ctx.PostForm("action_type")
	actionTypeValue, _ := strconv.Atoi(actionType)
	commentActionReq.ActionType = int32(actionTypeValue)

	// 创建评论操作
	if commentActionReq.ActionType == 1 {
		commentActionReq.CommentText = ctx.PostForm("comment_text")
	} else {
		commentId := ctx.PostForm("comment_id")
		commentActionReq.CommentId, _ = strconv.ParseInt(commentId, 10, 64)
	}

	videoServiceClient := ctx.Keys["video_service"].(service.VideoServiceClient)
	videoServiceResp, _ := videoServiceClient.CommentAction(context.Background(), &commentActionReq)

	// todo 这里太复杂了，找找能直接序列化不
	var comment res.Comment
	comment.ID = videoServiceResp.Comment.Id
	comment.CreateDate = videoServiceResp.Comment.CreateDate
	comment.Content = videoServiceResp.Comment.Content

	r := res.CommentActionResponse{
		StatusCode: videoServiceResp.StatusCode,
		StatusMsg:  videoServiceResp.StatusMsg,
		Comment:    &comment,
	}

	ctx.JSON(http.StatusOK, r)
}
