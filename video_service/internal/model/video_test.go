package model

import "testing"

func TestVideoModel_AddFavoriteCount(t *testing.T) {
	InitDb()
	GetVideoInstance().AddFavoriteCount(1949953677991936)
}

func TestFavoriteModel_DeleteFavorite2(t *testing.T) {
	InitDb()
	GetVideoInstance().DeleteFavoriteCount(1949953677991936)
}
