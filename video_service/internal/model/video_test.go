package model

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"
)

// 测试视频喜欢记录 + 1
func TestVideoModel_AddFavoriteCount(t *testing.T) {
	InitDb()
	tx := DB.Begin()
	GetVideoInstance().AddFavoriteCount(tx, 1949953677991936)
}

// 测试视频喜欢记录 - 1
func TestFavoriteModel_DeleteFavorite2(t *testing.T) {
	InitDb()
	tx := DB.Begin()
	GetVideoInstance().DeleteFavoriteCount(tx, 1949953677991936)
}

// 测试新增喜欢记录
func TestFavoriteModel_AddFavorite(t *testing.T) {
	InitDb()
	tx := DB.Begin()
	favorite := Favorite{
		UserId:  123,
		VideoId: 456,
	}
	GetFavoriteInstance().AddFavorite(tx, &favorite)
}

// 测试软删除喜欢记录
func TestFavoriteModel_DeleteFavorite(t *testing.T) {
	InitDb()
	tx := DB.Begin()
	favorite := Favorite{
		UserId:  123,
		VideoId: 456,
	}
	GetFavoriteInstance().DeleteFavorite(tx, &favorite)
}

// 测试创建评论
func TestCommentModel_CreateComment(t *testing.T) {
	InitDb()
	tx := DB.Begin()
	comment := Comment{
		UserId:  111,
		VideoId: 222,
		Content: "喜欢",
	}
	GetCommentInstance().CreateComment(tx, &comment)
	tx.Commit()
}

// 测试删除评论
func TestCommentModel_DeleteComment(t *testing.T) {
	InitDb()
	GetCommentInstance().DeleteComment(8361782507610112)

}

// 测试删除评论
func TestCommentModel_CommentList(t *testing.T) {
	InitDb()
	commentList, _ := GetCommentInstance().CommentList(4395719587667968)
	fmt.Print(commentList)
}

func TestTime(t *testing.T) {
	currentTime := time.Now()
	fmt.Println(currentTime)
	timeString := currentTime.Format("2006-01-02 15:04:05")
	fmt.Println("Formatted time:", timeString)
}

// 测试用户获赞数量
func TestVideoModel_GetFavoritedCount(t *testing.T) {
	InitDb()
	count, _ := GetVideoInstance().GetFavoritedCount(812575311663104)
	fmt.Println(count)
}

// 统计作品数量
func TestVideoModel_GetWorkCount(t *testing.T) {
	InitDb()
	count, _ := GetVideoInstance().GetWorkCount(812575311663104)
	fmt.Println(count)
}

// 统计喜欢数量
func TestFavorite_GetFavoriteCount(t *testing.T) {
	InitDb()
	count, _ := GetFavoriteInstance().GetFavoriteCount(812575311663104)
	fmt.Println(count)
}

// 找到喜欢的视频id
func TestFavoriteModel_FavoriteVideoList(t *testing.T) {
	InitDb()
	list, _ := GetFavoriteInstance().FavoriteVideoList(812575311663104)

	videoList, _ := GetVideoInstance().GetVideoList(list)
	fmt.Print(list)
	fmt.Print(videoList)
}

func TestFavoriteModel_IsFavorite(t *testing.T) {
	InitDb()
	favorite, _ := GetFavoriteInstance().IsFavorite(812575311663104, 2276964627783680)
	fmt.Print(favorite)
}

// 根据时间查找视频列表
func TestVideoModel_GetVideoByTime(t *testing.T) {
	InitDb()
	videos, _ := GetVideoInstance().GetVideoByTime(time.Now())
	fmt.Print(videos)
}

func TestFmt(t *testing.T) {
	userId := int64(111)
	key := fmt.Sprintf("%s:%s:%s", "user", "info", strconv.FormatInt(userId, 10))
	print(key)
}

func TestFavoriteModel_FavoriteUserList(t *testing.T) {
	InitDb()
	list, _ := GetFavoriteInstance().FavoriteUserList(4396360053694464)
	fmt.Println(list)
}

func TestCommentModel_GetComment(t *testing.T) {
	InitDb()
	tx := DB.Begin()
	comment, _ := commentModel.GetComment(tx, 4419719369990144)
	log.Print(comment)
}
