package model

import (
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
func (*CommentModel) CreateComment(comment *Comment) (id int64, err error) {
	flake, _ := snowFlake.NewSnowFlake(7, 2)
	comment.Id = flake.NextId()
	comment.CommentStatus = true
	comment.CreatAt = time.Now()

	result := DB.Create(&comment)
	if result.Error != nil {
		return -1, result.Error
	}

	return comment.Id, nil
}

// DeleteComment 删除评论
func (*CommentModel) DeleteComment(commentId int64) error {
	var comment Comment
	result := DB.First(&comment, commentId)
	if result.Error != nil {
		return result.Error
	}

	comment.CommentStatus = false
	result = DB.Save(&comment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
