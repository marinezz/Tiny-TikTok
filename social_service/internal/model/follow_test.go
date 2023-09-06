package model

import (
	"fmt"
	"testing"
)

func TestFollowModel_FollowAction(t *testing.T) {
	InitDb()
	follow := Follow{
		UserId:   1,
		ToUserId: 2,
		IsFollow: 2,
	}
	err := GetFollowInstance().FollowAction(&follow)
	fmt.Println(err)
}

func TestFollowModel_IsFollow(t *testing.T) {
	InitDb()
	res, err := GetFollowInstance().IsFollow(5155644223918080, 5155317378584576)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func TestFollowModel_GetFollowList(t *testing.T) {
	InitDb()
	var UserId []int64
	var reqUser int64
	reqUser = 1
	if err := GetFollowInstance().GetFollowList(reqUser, &UserId); err != nil {
		panic(nil)
	}
	fmt.Printf("%d的关注有%#v\n", reqUser, UserId)
}

func TestFollowModel_GetFollowerList(t *testing.T) {
	InitDb()
	var UserId []int64
	var reqUser int64
	reqUser = 2
	if err := GetFollowInstance().GetFollowerList(reqUser, &UserId); err != nil {
		panic(nil)
	}
	fmt.Printf("%d的粉丝有%#v\n", reqUser, UserId)
}

func TestFollowModel_GetFriendList(t *testing.T) {
	InitDb()
	var UserId []int64
	var reqUser int64
	reqUser = 2
	if err := GetFollowInstance().GetFriendList(reqUser, &UserId); err != nil {
		panic(nil)
	}
	fmt.Printf("%d的好友有%#v\n", reqUser, UserId)
}

func TestFollowModel_GetFollowCount(t *testing.T) {
	InitDb()
	var cnt int64
	var reqUser int64
	reqUser = 1
	cnt, err := GetFollowInstance().GetFollowCount(reqUser)
	if err != nil {
		panic(nil)
	}
	fmt.Printf("%d的关注有%d个\n", reqUser, cnt)
}

func TestFollowModel_GetFollowerCount(t *testing.T) {
	InitDb()
	var cnt int64
	var reqUser int64
	reqUser = 2
	cnt, err := GetFollowInstance().GetFollowerCount(reqUser)
	if err != nil {
		panic(nil)
	}
	fmt.Printf("%d的粉丝有%d个\n", reqUser, cnt)
}
