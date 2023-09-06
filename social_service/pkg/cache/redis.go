package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"strconv"
)

var Redis *redis.Client
var ctx = context.Background()

// InitRedis 连接redis
func InitRedis() {
	addr := viper.GetString("redis.address")
	Redis = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0, // 存入DB0
	})
	_, err := Redis.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}

func GenerateFollowKey(UserId int64) string {
	return "TINYTIKTOK:FOLLOW:" + strconv.FormatInt(UserId, 10)
}

func GenerateFollowerKey(UserId int64) string {
	return "TINYTIKTOK:FOLLOWER:" + strconv.FormatInt(UserId, 10)
}

func GenerateFriendKey(UserId int64) string {
	return "TINYTIKTOK:FRIEND:" + strconv.FormatInt(UserId, 10)
}

func IsFollow(UserId int64, ToUserId int64) (bool, error) {
	result, err := Redis.SIsMember(ctx, GenerateFollowKey(UserId), ToUserId).Result()
	return result, err
}

func FollowAction(UserId int64, ToUserId int64, ActionType int32) error {
	// 先查询关注的状态
	result, err := Redis.SIsMember(ctx, GenerateFollowKey(UserId), ToUserId).Result()
	if err != nil {
		return err
	}
	// 获取操作类型
	action := true
	if ActionType == 2 {
		action = false
	}
	// 判断操作是否需要执行
	if result == action {
		return nil
	}
	// 执行操作
	pipe := Redis.TxPipeline()
	defer pipe.Close()
	if action {
		pipe.SAdd(ctx, GenerateFollowKey(UserId), ToUserId)
		pipe.SAdd(ctx, GenerateFollowerKey(ToUserId), UserId)
		_, err := pipe.Exec(ctx)
		if err != nil {
			return err
		}
	} else {
		pipe.SRem(ctx, GenerateFollowKey(UserId), ToUserId)
		pipe.SRem(ctx, GenerateFollowerKey(ToUserId), UserId)
		_, err := pipe.Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetFollowList(reqUser int64, UserId *[]int64) error {
	result, err := Redis.SMembers(ctx, GenerateFollowKey(reqUser)).Result()
	if err != nil {
		return err
	}
	for _, r := range result {
		r64, _ := strconv.ParseInt(r, 10, 64)
		*UserId = append(*UserId, r64)
	}
	return nil
}

func GetFollowerList(reqUser int64, UserId *[]int64) error {
	result, err := Redis.SMembers(ctx, GenerateFollowerKey(reqUser)).Result()
	if err != nil {
		return err
	}
	for _, r := range result {
		r64, _ := strconv.ParseInt(r, 10, 64)
		*UserId = append(*UserId, r64)
	}
	return nil
}

func GetFriendList(reqUser int64, UserId *[]int64) error {
	_, err := Redis.Do(ctx, "SINTERSTORE", GenerateFriendKey(reqUser), GenerateFollowKey(reqUser), GenerateFollowerKey(reqUser)).Result()
	if err != nil {
		return err
	}
	result, err := Redis.SMembers(ctx, GenerateFriendKey(reqUser)).Result()
	if err != nil {
		return err
	}
	for _, r := range result {
		r64, _ := strconv.ParseInt(r, 10, 64)
		*UserId = append(*UserId, r64)
	}
	return nil
}

func GetFollowCount(reqUser int64) (int64, error) {
	return Redis.SCard(ctx, GenerateFollowKey(reqUser)).Result()
}

func GetFollowerCount(reqUser int64) (int64, error) {
	return Redis.SCard(ctx, GenerateFollowerKey(reqUser)).Result()
}
