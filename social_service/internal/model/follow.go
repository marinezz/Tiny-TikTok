package model

import (
	"errors"
	"gorm.io/gorm"
	"sync"
	"utils/snowFlake"
)

type Follow struct {
	Id       int64 `gorm:"primary_key"`
	UserId   int64
	ToUserId int64
	IsFollow int32 `gorm:"default:(2)"`
}

type FollowModel struct {
}

var followModel *FollowModel
var followOnce sync.Once // 单例模式

// GetFollowInstance 获取单例实例
func GetFollowInstance() *FollowModel {
	followOnce.Do(
		func() {
			followModel = &FollowModel{}
		},
	)
	return followModel
}

// FollowAction 更新关注状态
func (*FollowModel) FollowAction(follow *Follow) error {
	isFollow := follow.IsFollow
	if err := DB.Where(&Follow{UserId: follow.UserId, ToUserId: follow.ToUserId}).First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			flake, _ := snowFlake.NewSnowFlake(7, 3)
			follow.Id = flake.NextId()
			err = DB.Create(&follow).Error // create new record from newUser
		}
		if err != nil {
			return err
		}
	}
	if follow.Id != 0 && isFollow != follow.IsFollow {
		err := DB.Model(&Follow{}).Where(&follow).Update("is_follow", isFollow).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (*FollowModel) IsFollow(UserId int64, ToUserId int64) (bool, error) {
	follow := Follow{UserId: UserId, ToUserId: ToUserId}
	err := DB.Where(&follow).First(&follow).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if follow.IsFollow == 1 {
		return true, nil
	}
	return false, nil
}

func (*FollowModel) GetFollowList(reqUser int64, UserId *[]int64) error {
	var follows []Follow
	if err := DB.Where(&Follow{UserId: reqUser, IsFollow: 1}).Find(&follows).Error; err != nil {
		return err
	}

	for _, follow := range follows {
		*UserId = append(*UserId, follow.ToUserId)
	}
	return nil
}

func (*FollowModel) GetFollowerList(reqUser int64, UserId *[]int64) error {
	var follows []Follow
	if err := DB.Where(&Follow{ToUserId: reqUser, IsFollow: 1}).Find(&follows).Error; err != nil {
		return err
	}

	for _, follow := range follows {
		*UserId = append(*UserId, follow.UserId)
	}
	return nil
}

func (*FollowModel) GetFriendList(reqUser int64, UserId *[]int64) error {

	if err := DB.Raw("select a.to_user_id from follow as a inner join follow as b on a.user_id = b.to_user_id and a.to_user_id = b.user_id and a.is_follow = 1 and b.is_follow = 1 and a.user_id = ?", reqUser).Scan(UserId).Error; err != nil {
		return err
	}

	return nil
}

func (*FollowModel) GetFollowCount(reqUser int64) (int64, error) {
	var cnt int64
	if err := DB.Model(&Follow{}).Where(&Follow{UserId: reqUser, IsFollow: 1}).Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}

func (*FollowModel) GetFollowerCount(reqUser int64) (int64, error) {
	var cnt int64
	if err := DB.Model(&Follow{}).Where(&Follow{ToUserId: reqUser, IsFollow: 1}).Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}
