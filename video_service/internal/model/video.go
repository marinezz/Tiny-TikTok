package model

import "sync"

type Video struct {
	Id            int64 `gorm:"primary_key"`
	AuthId        string
	Title         string
	CoverUrl      string `gorm:"default:(-)"`
	PlayUrl       string `gorm:"default:(-)"`
	FavoriteCount int    `gorm:"default:(-)"`
	CommentCount  int    `gorm:"default:(-)"`
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
