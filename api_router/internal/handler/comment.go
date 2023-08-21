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
	videoServiceResp, err := videoServiceClient.CommentAction(context.Background(), &commentActionReq)
	if err != nil {
		PanicIfCommentError(err)
	}

	if actionTypeValue == 1 {
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
	// 如果是删除评论的操作
	if actionTypeValue == 2 {
		r := res.CommentDeleteResponse{
			StatusCode: videoServiceResp.StatusCode,
			StatusMsg:  videoServiceResp.StatusMsg,
		}

		ctx.JSON(http.StatusOK, r)
	}
}

func CommentList(ctx *gin.Context) {
	var commentListReq service.CommentListRequest

	videoIdStr := ctx.Query("video_id")
	videoId, _ := strconv.ParseInt(videoIdStr, 10, 64)

	commentListReq.VideoId = videoId

	videoServiceClient := ctx.Keys["video_service"].(service.VideoServiceClient)
	commentListResp, err := videoServiceClient.CommentList(context.Background(), &commentListReq)
	if err != nil {
		PanicIfCommentError(err)
	}

	// 找到所有的用户Id
	var userIds []int64
	for _, comment := range commentListResp.CommentList {
		userIds = append(userIds, comment.UserId)
	}

	userInfos := GetUserInfo(userIds, ctx)

	commentList := BuildCommentList(commentListResp.CommentList, userInfos)

	r := res.CommentListResponse{
		StatusCode: commentListResp.StatusCode,
		StatusMsg:  commentListResp.StatusMsg,
		Comments:   commentList,
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

func BuildCommentList(comments []*service.Comment, userInfos []res.User) []res.Comment {
	var commentList []res.Comment

	for i, comment := range comments {
		commentList = append(commentList, res.Comment{
			Id:         comment.Id,
			User:       userInfos[i],
			Content:    comment.Content,
			CreateDate: comment.CreateDate,
		})
	}

	return commentList
}
