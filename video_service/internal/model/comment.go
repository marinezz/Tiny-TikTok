package model

import (
	"sync"
	"time"
)

type Comment struct {
	Id            int64 `gorm:"primaryKey"`
	UserId        int64
	VideoId       int64
	content       string `gorm:"default:(-)"` // 评论内容
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
