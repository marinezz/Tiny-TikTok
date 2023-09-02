package handler

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
	"video/internal/model"
	"video/pkg/cache"
)

func TestVideoService_Feed(t *testing.T) {
	timePoint := time.Unix(1693130905305/1000, 0)
	fmt.Print(timePoint)
}

func TestCache(t *testing.T) {
	model.InitDb()
	cache.InitRedis()
	key := "video:comment_list:4395719587667968"

	comment, _ := model.GetCommentInstance().GetComment(5146154783088640)
	commentJson, _ := json.Marshal(comment)

	removedCount, err := cache.Redis.ZRem(cache.Ctx, key, string(commentJson)).Result()
	if err != nil {
		fmt.Println("Redis操作出错：", err)
	} else {
		fmt.Printf("删除了 %d 个匹配的成员\n", removedCount)
	}
}
