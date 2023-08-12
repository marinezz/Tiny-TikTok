package model

import (
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
	FavoriteCount int    `gorm:"default:(-)"`
	CommentCount  int    `gorm:"default:(-)"`
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
