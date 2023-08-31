package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"utils/exception"
	"video/internal/model"
	"video/internal/service"
	"video/pkg/cache"
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

		tx := model.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// 事务中执行创建操作
		id, err := model.GetCommentInstance().CreateComment(tx, &comment)
		if err != nil {
			tx.Rollback()
			resp.StatusCode = exception.SUCCESS
			resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
			resp.Comment = nil
			return resp, err
		}

		// 视频评论数量 + 1
		err = model.GetVideoInstance().AddCommentCount(tx, req.VideoId)
		if err != nil {
			tx.Rollback()
			resp.StatusCode = exception.SUCCESS
			resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
			resp.Comment = nil
			return resp, err
		}

		tx.Commit()

		commentResp := &service.Comment{
			Id:      id,
			Content: req.CommentText,
			// 将Time.time转换成字符串形式，格式为mm-dd
			CreateDate: time.Format("01-02"),
		}

		// 将评论返回
		resp.StatusCode = exception.SUCCESS
		resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
		resp.Comment = commentResp

		return resp, nil
	}

	// 删除评论
	comment.CreatAt = time

	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = model.GetCommentInstance().DeleteComment(tx, req.CommentId)
	if err != nil {
		resp.StatusCode = exception.CommentDeleteErr
		resp.StatusMsg = exception.GetMsg(exception.CommentDeleteErr)
		resp.Comment = nil

		return resp, err
	}
	// 视频评论数量 - 1
	err = model.GetVideoInstance().DeleteCommentCount(tx, req.VideoId)
	if err != nil {
		resp.StatusCode = exception.CommentDeleteErr
		resp.StatusMsg = exception.GetMsg(exception.CommentDeleteErr)
		resp.Comment = nil

		return resp, err
	}
	tx.Commit()

	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	resp.Comment = nil

	return resp, nil
}

// CommentList 评论列表
func (*VideoService) CommentList(ctx context.Context, req *service.CommentListRequest) (resp *service.CommentListResponse, err error) {
	resp = new(service.CommentListResponse)
	var comments []model.Comment
	key := fmt.Sprintf("%s:%s:%s", "video", "comment_list", req.VideoId)

	exist, err := cache.Redis.Exists(context.Background(), key).Result()
	if err != nil {
		resp.StatusCode = exception.CacheErr
		resp.StatusMsg = exception.GetMsg(exception.CacheErr)
		return resp, err
	}

	if exist > 0 {
		commentsString, err := cache.Redis.Get(context.Background(), key).Result()
		if err != nil {
			resp.StatusCode = exception.VideoUnExist
			resp.StatusMsg = exception.GetMsg(exception.VideoUnExist)
			return resp, err
		}
		err = json.Unmarshal([]byte(commentsString), &comments)
		if err != nil {
			return nil, err
		}
	} else {
		// 根据视频id找到所有的评论
		comments, err = model.GetCommentInstance().CommentList(req.VideoId)
		if err != nil {
			resp.StatusCode = exception.CommentUnExist
			resp.StatusMsg = exception.GetMsg(exception.CommentUnExist)
			return nil, err
		}

		// 将查询结果放入缓存中
		commentJson, _ := json.Marshal(&comments)
		err = cache.Redis.Set(context.Background(), key, commentJson, 30*time.Minute).Err()
		if err != nil {
			resp.StatusCode = exception.CacheErr
			resp.StatusMsg = exception.GetMsg(exception.CacheErr)
			return resp, err
		}
	}

	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	resp.CommentList = BuildComments(comments)

	return resp, nil
}

func BuildComments(comments []model.Comment) []*service.Comment {
	var commentresp []*service.Comment

	for _, comment := range comments {
		commentresp = append(commentresp, &service.Comment{
			Id:         comment.Id,
			UserId:     comment.UserId,
			Content:    comment.Content,
			CreateDate: comment.CreatAt.Format("01-02"),
		})
	}

	return commentresp
}
