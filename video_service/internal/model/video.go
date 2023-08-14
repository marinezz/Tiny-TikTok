package model

import (
	"gorm.io/gorm"
	"sync"
	"time"
	"utils/snowFlake"
)

type Video struct {
	Id            int64 `gorm:"primary_key"`
	AuthId        int64
	Title         string
	CoverUrl      string `gorm:"default:(-)"`
	PlayUrl       string `gorm:"default:(-)"`
	FavoriteCount int    `gorm:"default:0"`
	CommentCount  int    `gorm:"default:0"`
	CreatAt       time.Time
}

type VideoModel struct {
}

var videoMedel *VideoModel
var videoOnce sync.Once // 单例模式

// GetVideoInstance 获取单例的实例
func GetVideoInstance() *VideoModel {
	videoOnce.Do(
		func() {
			videoMedel = &VideoModel{}
		})
	return videoMedel
}

// Create 创建视频信息
func (*VideoModel) Create(video *Video) error {
	// 服务2
	flake, _ := snowFlake.NewSnowFlake(7, 2)
	video.Id = flake.NextId()
	DB.Create(&video)
	return nil
}

// AddFavoriteCount 喜欢记录 + 1
func (*VideoModel) AddFavoriteCount(videoId int64) error {
	result := DB.Model(&Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteFavoriteCount 喜欢记录 - 1
func (*VideoModel) DeleteFavoriteCount(videoId int64) error {
	result := DB.Model(&Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return nil
	return nil
}

// AddCommentCount 视频评论数量 + 1
func (*VideoModel) AddCommentCount(videoId int64) error {
	result := DB.Model(&Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteCommentCount 视频评论数量 - 1
func (*VideoModel) DeleteCommentCount(videoId int64) error {
	result := DB.Model(&Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return nil
}
