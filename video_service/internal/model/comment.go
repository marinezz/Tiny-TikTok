package model

import (
	"errors"
	"gorm.io/gorm"
	"sync"
	"time"
	"utils/snowFlake"
)

type Comment struct {
	Id            int64 `gorm:"primaryKey"`
	UserId        int64
	VideoId       int64
	Content       string `gorm:"default:(-)"` // 评论内容
	CreatAt       time.Time
	CommentStatus bool `gorm:"default:(-)"`
}

type CommentModel struct {
}

var commentModel *CommentModel
var commentOnce sync.Once

// GetCommentInstance 拿到单例实例
func GetCommentInstance() *CommentModel {
	commentOnce.Do(
		func() {
			commentModel = &CommentModel{}
		})
	return commentModel
}

// CreateComment 新增评论
func (*CommentModel) CreateComment(tx *gorm.DB, comment *Comment) (id int64, err error) {
	flake, _ := snowFlake.NewSnowFlake(7, 2)
	comment.Id = flake.NextId()
	comment.CommentStatus = true
	comment.CreatAt = time.Now()

	result := tx.Create(&comment)
	if result.Error != nil {
		return -1, result.Error
	}

	return comment.Id, nil
}

// DeleteComment 删除评论
func (*CommentModel) DeleteComment(tx *gorm.DB, commentId int64) error {
	var comment Comment
	result := tx.First(&comment, commentId)
	if result.Error != nil {
		return result.Error
	}
	if comment.CommentStatus == false {
		return errors.New("评论不存在！！")
	}

	comment.CommentStatus = false
	result = tx.Save(&comment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CommentList 根据视频id找到所有的评论
func (*CommentModel) CommentList(videoId int64) ([]Comment, error) {
	var comments []Comment

	result := DB.Table("comment").
		Where("video_id = ? AND comment_status = ?", videoId, true).
		Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}

	return comments, nil
}
