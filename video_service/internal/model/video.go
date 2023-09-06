package model

import (
	"gorm.io/gorm"
	"sync"
	"time"
	"utils/snowFlake"
)

type Video struct {
	Id       int64 `gorm:"primary_key"`
	AuthId   int64
	Title    string
	CoverUrl string `gorm:"default:(-)"`
	PlayUrl  string `gorm:"default:(-)"`
	//FavoriteCount int64  `gorm:"default:0"`
	//CommentCount int64 `gorm:"default:0"`
	CreatAt time.Time
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

// DeleteVideoByUrl 删除视频
func (v *VideoModel) DeleteVideoByUrl(videoUrl string) error {
	var video Video
	if err := DB.Where("play_url = ?", videoUrl).First(&video).Error; err != nil {
		return err
	}

	// 删除找到的记录
	if err := DB.Delete(&video).Error; err != nil {
		return err
	}

	return nil
}

// GetVideoByTime 根据创建时间获取视频
func (*VideoModel) GetVideoByTime(timePoint time.Time) ([]Video, error) {
	var videos []Video

	result := DB.Table("video").
		Where("creat_at < ?", timePoint).
		Order("creat_at DESC").
		Limit(30).
		Find(&videos)
	if result.Error != nil {
		return nil, result.Error
	}

	// 查询不到数据，就返回当前时间最新的30条数据
	if len(videos) == 0 {
		timePoint = time.Now()
		result := DB.Table("video").
			Where("creat_at < ?", timePoint).
			Order("creat_at DESC").
			Limit(30).
			Find(&videos)
		if result.Error != nil {
			return nil, result.Error
		}
		return videos, nil
	}

	return videos, nil
}

// GetVideoList 根据视频Id获取视频列表
func (*VideoModel) GetVideoList(videoIds []int64) ([]Video, error) {
	var videos []Video

	result := DB.Table("video").
		Where("id IN ?", videoIds).
		Find(&videos)
	if result.Error != nil {
		return nil, result.Error
	}

	return videos, nil
}

// GetVideoListByUser 根据用户的id找到视频列表
func (*VideoModel) GetVideoListByUser(userId int64) ([]Video, error) {
	var videos []Video

	result := DB.Table("video").
		Where("auth_id = ?", userId).
		Find(&videos)
	if result.Error != nil {
		return nil, result.Error
	}

	return videos, nil
}

// AddFavoriteCount 喜欢记录 + 1
func (*VideoModel) AddFavoriteCount(tx *gorm.DB, videoId int64) error {
	result := tx.Model(&Video{}).Where("id = ?", videoId).
		Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteFavoriteCount 喜欢记录 - 1
func (*VideoModel) DeleteFavoriteCount(tx *gorm.DB, videoId int64) error {
	result := tx.Model(&Video{}).Where("id = ?", videoId).
		Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// AddCommentCount 视频评论数量 + 1
func (*VideoModel) AddCommentCount(tx *gorm.DB, videoId int64) error {
	result := tx.Model(&Video{}).
		Where("id = ?", videoId).
		Update("comment_count", gorm.Expr("comment_count + ?", 1))
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteCommentCount 视频评论数量 - 1
func (*VideoModel) DeleteCommentCount(tx *gorm.DB, videoId int64) error {
	result := tx.Model(&Video{}).
		Where("id = ?", videoId).
		Update("comment_count", gorm.Expr("comment_count - ?", 1))
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetFavoritedCount 获取用户的获赞数量
func (*VideoModel) GetFavoritedCount(userId int64) (int64, error) {
	var count int64

	DB.Table("video").
		Where("auth_id=?", userId).
		Select("SUM(favorite_count) as count").
		Pluck("count", &count)

	return count, nil
}

// GetWorkCount 获取用户的作品数量
func (*VideoModel) GetWorkCount(userId int64) (int64, error) {
	var count int64
	DB.Table("video").
		Where("auth_id=?", userId).
		Count(&count)

	return count, nil
}
