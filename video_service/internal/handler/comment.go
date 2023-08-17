package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
	"utils/exception"
	"video/internal/model"
	"video/internal/service"
)

// CommentAction 评论操作
func (*VideoService) CommentAction(ctx context.Context, req *service.CommentActionRequest) (resp *service.CommentActionResponse, err error) {
	resp = new(service.CommentActionResponse)
	comment := model.Comment{
		UserId:  req.UserId,
		VideoId: req.VideoId,
		Content: req.CommentText,
	}
	action := req.ActionType

	time := time.Now()

	// 发布评论
	if action == 1 {
		comment.CreatAt = time
		id, _ := model.GetCommentInstance().CreateComment(&comment)

		// 视频评论数量 + 1
		model.GetVideoInstance().AddCommentCount(req.VideoId)

		commentResp := &service.Comment{
			Id:      id,
			Content: req.CommentText,
			// 将Time.time转换成字符串形式
			CreateDate: time.Format("2006-01-02 15:04:05"),
		}

		// 将评论返回
		resp.StatusCode = exception.SUCCESS
		resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
		resp.Comment = commentResp

		return resp, nil
	}

	// 删除评论
	model.GetCommentInstance().DeleteComment(req.CommentId)
	// 视频评论数量 - 1
	model.GetVideoInstance().DeleteCommentCount(req.VideoId)

	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	resp.Comment = nil

	return resp, nil
}

// CommentList 评论列表
func (*VideoService) CommentList(ctx context.Context, req *service.CommentListRequest) (resp *service.CommentListResponse, err error) {

	return nil, status.Errorf(codes.Unimplemented, "method CommentList not implemented")
}
