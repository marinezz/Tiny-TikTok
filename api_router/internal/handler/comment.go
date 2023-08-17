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
	commentActionReq.ActionType = int64(actionTypeValue)

	// 评论操作
	if commentActionReq.ActionType == 1 {
		commentActionReq.CommentText = ctx.PostForm("comment_text")
	} else {
		commentId := ctx.PostForm("comment_id")
		commentActionReq.CommentId, _ = strconv.ParseInt(commentId, 10, 64)
	}

	videoServiceClient := ctx.Keys["video_service"].(service.VideoServiceClient)
	videoServiceResp, _ := videoServiceClient.CommentAction(context.Background(), &commentActionReq)

	// 如果是删除评论的操作
	if actionTypeValue == 2 {
		r := res.CommentActionResponse{
			StatusCode: videoServiceResp.StatusCode,
			StatusMsg:  videoServiceResp.StatusMsg,
			//Comment:    nil,
		}

		ctx.JSON(http.StatusOK, r)
	}

	// 构建用户信息
	userIds := []int64{userId.(int64)}
	userInfos := GetUserInfo(userIds, ctx)

	r := res.CommentActionResponse{
		StatusCode: videoServiceResp.StatusCode,
		StatusMsg:  videoServiceResp.StatusMsg,
		Comment:    BuildComment(videoServiceResp.Comment, userInfos[0]),
	}

	ctx.JSON(http.StatusOK, r)
}

func BuildComment(comment *service.Comment, userInfo res.User) res.Comment {

	return res.Comment{
		Id:         comment.Id,
		User:       userInfo,
		Content:    comment.Content,
		CreateDate: comment.CreateDate,
	}
}
