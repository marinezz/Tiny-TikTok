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
	IsFavorite bool `gorm:"default:(-)"`
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
func (*FavoriteModel) AddFavorite(favorite *Favorite) error {
	result := DB.Where("user_id=? AND video_id=?", favorite.UserId, favorite.VideoId).First(&favorite)
	// 发生除没找到记录的其它错误
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	// 如果找到了记录，更新is_favorite置为1
	if result.RowsAffected > 0 {
		favorite.IsFavorite = true
		result = DB.Save(&favorite)
		if result.Error != nil {
			return result.Error
		}
	} else {
		// 否则创建新记录
		flake, _ := snowFlake.NewSnowFlake(7, 2)
		favorite.Id = flake.NextId()
		result = DB.Create(&favorite)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

// DeleteFavorite 删除点赞
func (*FavoriteModel) DeleteFavorite(favorite *Favorite) error {
	result := DB.Where("user_id=? AND video_id=?", favorite.UserId, favorite.VideoId).First(&favorite)
	// 发生除没找到记录的其它错误
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	// 如果找到了记录，更新is_favorite置为0
	if result.RowsAffected > 0 {
		favorite.IsFavorite = false
		result = DB.Save(&favorite)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
