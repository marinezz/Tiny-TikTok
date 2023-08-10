package model

import "sync"

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
