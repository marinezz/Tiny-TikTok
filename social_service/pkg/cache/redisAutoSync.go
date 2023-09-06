package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"social/internal/model"
	"time"
)

// getAllKeys 获取所有"FOLLOW"表项的key值
func getAllKeys(keys *[]string) {
	// 使用正则表达式来匹配键
	pattern := "TINYTIKTOK:FOLLOW:*" // 替换为您的正则表达式

	// 初始化SCAN迭代器
	iter := Redis.Scan(ctx, 0, pattern, 0).Iterator()

	// 用于并发处理的函数
	handleKey := func(key string) error {
		*keys = append(*keys, key)
		return nil
	}

	// 并发处理匹配的键
	var g errgroup.Group
	for iter.Next(ctx) {
		key := iter.Val()
		g.Go(func() error {
			return handleKey(key)
		})
	}

	if err := g.Wait(); err != nil {
		panic(err)
	}
}

func getAllValueByKeys(keys []string) string {
	// 创建管道并将所有的 查询列表操作加入pipe
	cmds, err := Redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			pipe.SMembers(ctx, key)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	followInfo := make(map[string][]string)
	for index, cmd := range cmds {
		followInfo[keys[index][18:]] = cmd.(*redis.StringSliceCmd).Val()
	}
	marshal, err := json.Marshal(followInfo)
	if err != nil {
		return ""
	}
	return string(marshal)

}

/*
	func MysqlToRedis() error {
		var errs []error
		// 所有数据，每次取100条
		page := 1
		pageSize := 100
		for {
			var follows []Follow
			offset := (page - 1) * pageSize
			DB.Limit(pageSize).Offset(offset).Find(&follows)

			// 处理查询结果
			if len(follows) == 0 {
				break // 没有更多记录可获取，退出循环
			}

			for _, follow := range follows {
				err := cache.FollowAction(follow.UserId, follow.ToUserId, 1)
				errs = append(errs, err)
			}

			page++
		}

		log.Println(errs)
		return nil
*/
func AutoSync() {
	// 先匹配所有follow的key
	var keys []string
	getAllKeys(&keys)
	res := getAllValueByKeys(keys)
	err := model.GetFollowInstance().RedisToMysql(res)
	if err != nil {
		panic(err)
	}
}

func TimerSync() {
	c := cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)), cron.WithLogger(
		cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))

	_, err := c.AddFunc("0 */10 * * * *", func() {
		fmt.Println(time.Now(), " 开始同步到mysql数据库")
		AutoSync()
	})
	if err != nil {
		panic(err)
	}

	c.Start()
	select {}
}
