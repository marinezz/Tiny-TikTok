package cache

import (
	"fmt"
	"social/internal/model"
	"testing"
)

func TestFollowAction(t *testing.T) {
	InitRedis()
	err := FollowAction(3, 4, 1)
	fmt.Println(err)
}

func TestGetFriendListList(t *testing.T) {
	InitRedis()
	var UserId []int64
	err := GetFollowList(1, &UserId)
	fmt.Printf("%#v", UserId)
	fmt.Println(err)
}

func TestGetFollowCount(t *testing.T) {
	InitRedis()
	fmt.Println(GetFollowCount(100))
}

func TestAutoSync(t *testing.T) {
	InitRedis()
	model.InitDb()
	AutoSync()
}
