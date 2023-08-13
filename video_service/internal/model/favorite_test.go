package model

import "testing"

func TestFavoriteModel_AddFavorite(t *testing.T) {
	InitDb()
	favorite := Favorite{
		UserId:  123,
		VideoId: 456,
	}
	GetFavoriteInstance().AddFavorite(&favorite)
}

func TestFavoriteModel_DeleteFavorite(t *testing.T) {
	InitDb()
	favorite := Favorite{
		UserId:  123,
		VideoId: 456,
	}
	GetFavoriteInstance().DeleteFavorite(&favorite)
}
