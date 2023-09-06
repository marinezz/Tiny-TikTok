package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
	"utils/exception"
	"video/internal/model"
	"video/internal/service"
	"video/pkg/cache"
)

// CommentAction 评论操作
func (*VideoService) CommentAction(ctx context.Context, req *service.CommentActionRequest) (resp *service.CommentActionResponse, err error) {
	resp = new(service.CommentActionResponse)
	key := fmt.Sprintf("%s:%s:%s", "video", "comment_list", strconv.FormatInt(req.VideoId, 10))
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

		// 事务中执行创建操作
		tx := model.DB.Begin()
		id, err := model.GetCommentInstance().CreateComment(tx, &comment)
		if err != nil {
			resp.StatusCode = exception.CommentErr
			resp.StatusMsg = exception.GetMsg(exception.CommentErr)
			resp.Comment = nil
			return resp, err
		}
		comment.Id = id
		comment.CommentStatus = true
		commentJson, _ := json.Marshal(comment)

		// 存入缓存中
		member := redis.Z{
			Score:  float64(time.Unix()),
			Member: commentJson,
		}

		err = cache.Redis.ZAdd(cache.Ctx, key, &member).Err()
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("缓存错误：%v", err)
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
	tx := model.DB.Begin()
	commentInstance, err := model.GetCommentInstance().GetComment(tx, req.CommentId)
	commentMarshal, _ := json.Marshal(commentInstance)

	err = model.GetCommentInstance().DeleteComment(req.CommentId)
	if err != nil {
		resp.StatusCode = exception.CommentDeleteErr
		resp.StatusMsg = exception.GetMsg(exception.CommentDeleteErr)
		resp.Comment = nil

		return resp, err
	}
	log.Print(commentMarshal)

	// 删除缓存
	count, err := cache.Redis.ZRem(cache.Ctx, key, string(commentMarshal)).Result()
	log.Printf("删除了： %v 条记录", count)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("缓存错误：%v", err)
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
	key := fmt.Sprintf("%s:%s:%s", "video", "comment_list", strconv.FormatInt(req.VideoId, 10))

	exist, err := cache.Redis.Exists(cache.Ctx, key).Result()
	if err != nil {
		log.Print(4)
		return nil, fmt.Errorf("缓存错误：%v", err)
	}

	if exist == 0 {
		err := buildCommentCache(req.VideoId)
		if err != nil {
			log.Print(3)
			return nil, fmt.Errorf("缓存错误：%v", err)
		}
	}

	// 查询缓存
	commentsString, err := cache.Redis.ZRevRange(cache.Ctx, key, 0, -1).Result()
	if err != nil {
		log.Print(5)
		return nil, fmt.Errorf("缓存错误：%v", err)
	}

	for _, commentString := range commentsString {
		var comment model.Comment
		err := json.Unmarshal([]byte(commentString), &comment)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
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

// 构建评论列表缓存
func buildCommentCache(videoId int64) error {
	key := fmt.Sprintf("%s:%s:%s", "video", "comment_list", strconv.FormatInt(videoId, 10))

	comments, err := model.GetCommentInstance().CommentList(videoId)
	if err != nil {
		return err
	}

	var zMembers []*redis.Z
	for _, comment := range comments {
		commentJSON, err := json.Marshal(comment)
		if err != nil {
			fmt.Println("Error encoding comment:", err)
			continue
		}
		zMembers = append(zMembers, &redis.Z{
			Score:  float64(comment.CreatAt.Unix()),
			Member: commentJSON,
		})
	}

	err = cache.Redis.ZAdd(cache.Ctx, key, zMembers...).Err()
	if err != nil {
		return err
	}

	return nil
}

// 通过缓存查看视频得评论数量
func getCommentCount(videoId int64) int64 {
	key := fmt.Sprintf("%s:%s:%s", "video", "comment_list", strconv.FormatInt(videoId, 10))

	exists, err := cache.Redis.Exists(cache.Ctx, key).Result()
	if err != nil {
		log.Print(err)
	}

	if exists == 0 {
		err := buildCommentCache(videoId)
		if err != nil {
			log.Print(err)
		}
	}

	count, err := cache.Redis.ZCard(cache.Ctx, key).Result()
	if err != nil {
		log.Print(err)
	}

	return count
}
