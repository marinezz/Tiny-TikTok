package model

import (
	"errors"
	"gorm.io/gorm"
	"sync"
	"utils/snowFlake"
)

type Favorite struct {
	Id         int64 `gorm:"primaryKey"`
	UserId     int64
	VideoId    int64
	IsFavorite bool `gorm:"default:true"`
}

type FavoriteModel struct {
}

var favoriteModel *FavoriteModel
var favoriteOnce sync.Once

func GetFavoriteInstance() *FavoriteModel {
	favoriteOnce.Do(
		func() {
			favoriteModel = &FavoriteModel{}
		})
	return favoriteModel
}

// AddFavorite 创建点赞
func (*FavoriteModel) AddFavorite(favorite *Favorite) (bool, error) {
	result := DB.Where("user_id=? AND video_id=?", favorite.UserId, favorite.VideoId).First(&favorite)
	// 发生除没找到记录的其它错误
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, result.Error
	}

	// 判断是否需要视频表中创建新的记录
	isAdd := false

	// 如果找到了记录，更新is_favorite置为1
	if result.RowsAffected > 0 {
		if favorite.IsFavorite == false {
			isAdd = true
			favorite.IsFavorite = true
			result = DB.Save(&favorite)
			if result.Error != nil {
				return false, result.Error
			}
		}
	} else {
		// 否则创建新记录
		flake, _ := snowFlake.NewSnowFlake(7, 2)
		favorite.Id = flake.NextId()
		result = DB.Create(&favorite)
		isAdd = true
		if result.Error != nil {
			return false, result.Error
		}
	}
	return isAdd, nil
}

// IsFavorite 根据用户id和视频id获取点赞状态
func (*FavoriteModel) IsFavorite(userId int64, videoId int64) (bool, error) {
	var isFavorite bool

	result := DB.Table("favorite").
		Where("user_id = ? AND video_id = ?", userId, videoId).
		Pluck("is_favorite", &isFavorite)
	if result.Error != nil {
		return true, result.Error
	}

	return isFavorite, nil
}

// DeleteFavorite 删除点赞
func (*FavoriteModel) DeleteFavorite(favorite *Favorite) (error, bool) {
	result := DB.Where("user_id=? AND video_id=?", favorite.UserId, favorite.VideoId).First(&favorite)
	// 发生除没找到记录的其它错误
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error, false
	}
	isDelete := false
	// 如果找到了记录，更新is_favorite置为0
	if result.RowsAffected > 0 {
		if favorite.IsFavorite == true {
			isDelete = true

			favorite.IsFavorite = false
			result = DB.Save(&favorite)
			if result.Error != nil {
				return result.Error, false
			}
		}
	}

	return nil, isDelete
}

// FavoriteVideoList 根据用户Id获取所有喜欢的视频id
func (*FavoriteModel) FavoriteVideoList(userId int64) ([]int64, error) {
	var videoIds []int64

	result := DB.Table("favorite").
		Where("user_id = ? AND is_favorite = ?", userId, true).
		Pluck("video_id", &videoIds)
	if result.Error != nil {
		return nil, result.Error
	}

	return videoIds, nil
}

// GetFavoriteCount 获取喜欢数量
func (*FavoriteModel) GetFavoriteCount(userId int64) (int64, error) {
	var count int64

	DB.Table("favorite").
		Where("user_id=? AND is_favorite=?", userId, true).
		Count(&count)

	return count, nil
}
