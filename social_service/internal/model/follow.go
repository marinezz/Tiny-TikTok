package model

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"log"
	"strconv"
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

func (*FollowModel) RedisToMysql(data string) error {
	// 先将所有follow置为2
	if err := DB.Model(&Follow{}).Where("is_follow = ?", 1).Update("is_follow", 2).Error; err != nil {
		return err
	}
	// 转换data为map
	var follows map[string][]string
	if err := json.Unmarshal([]byte(data), &follows); err != nil {
		return err
	}
	// 对每一个关注进行操作
	var wg sync.WaitGroup
	var errs []error

	for userId, toUserIds := range follows {
		userid, _ := strconv.ParseInt(userId, 10, 64)
		for _, toUserId := range toUserIds {
			touserid, _ := strconv.ParseInt(toUserId, 10, 64)
			wg.Add(1)
			go func(userId int64, toUserId int64) {
				defer wg.Done()
				follow := Follow{UserId: userId, ToUserId: toUserId, IsFollow: 1}
				err := GetFollowInstance().FollowAction(&follow)
				if err != nil {
					errs = append(errs, err)
				}
			}(userid, touserid)

		}
	}

	wg.Wait()

	log.Println(errs)

	return nil
}
