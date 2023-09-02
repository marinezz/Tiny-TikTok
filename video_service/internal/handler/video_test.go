package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
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

func TestFavoriteCache(t *testing.T) {
	model.InitDb()
	cache.InitRedis()

	key := "video:favorite_video:4410007823982592"

	_, err := cache.Redis.Del(cache.Ctx, key).Result()
	if err != nil {
		fmt.Println("Error deleting existing Set:", err)
		return
	}

	_, err = cache.Redis.SAdd(cache.Ctx, key, "tempMember").Result()
	if err != nil {
		fmt.Println("Error adding member to Set:", err)
		return
	}

	_, err = cache.Redis.Expire(cache.Ctx, key, 12*time.Hour).Result()
	if err != nil {
		fmt.Println("Error setting expiration time:", err)
		return
	}

	// 检查 Set 是否存在
	exists, err := cache.Redis.Exists(cache.Ctx, key).Result()
	if err != nil {
		fmt.Println("Error checking Set existence:", err)
		return
	}

	if exists == 1 {
		fmt.Println("Set exists")
	} else {
		fmt.Println("Set does not exist")
	}

	// 删除临时成员并检查是否删除成功
	removedCount, err := cache.Redis.SRem(cache.Ctx, key, "tempMember").Result()
	if err != nil {
		fmt.Println("Error removing temporary member:", err)
		return
	}

	if removedCount == 1 {
		fmt.Println("Successfully removed temporary member")
	} else {
		fmt.Println("Temporary member not found in Set")
		return
	}

	// 检查 Set 是否存在
	exists, err = cache.Redis.Exists(cache.Ctx, key).Result()
	if err != nil {
		fmt.Println("Error checking Set existence:", err)
		return
	}

	if exists == 1 {
		fmt.Println("Set exists")
	} else {
		fmt.Println("Set does not exist")
	}
}

func TestSet(t *testing.T) {
	model.InitDb()
	cache.InitRedis()

	key := "video:favorite_video:4410007823982592"

	// 先添加临时成员
	_, err := cache.Redis.SAdd(cache.Ctx, key, "tempMember").Result()
	if err != nil {
		fmt.Println("Error adding member to Set:", err)
		return
	}

	// 检查 Set 是否存在
	exists, err := cache.Redis.Exists(cache.Ctx, key).Result()
	if err != nil {
		fmt.Println("Error checking Set existence:", err)
		return
	}

	if exists == 1 {
		fmt.Println("Set exists")
	} else {
		fmt.Println("Set does not exist")
		return
	}

	// 删除临时成员并检查是否删除成功
	removedCount, err := cache.Redis.SRem(cache.Ctx, key, "tempMember").Result()
	if err != nil {
		fmt.Println("Error removing temporary member:", err)
		return
	}

	if removedCount == 1 {
		fmt.Println("Successfully removed temporary member")
	} else {
		fmt.Println("Temporary member not found in Set")
		return
	}

	// 再次检查 Set 是否存在
	exists, err = cache.Redis.Exists(cache.Ctx, key).Result()
	if err != nil {
		fmt.Println("Error checking Set existence:", err)
		return
	}

	if exists == 1 {
		fmt.Println("Set exists")
	} else {
		fmt.Println("Set does not exist")
	}

}

func TestZset(t *testing.T) {
	model.InitDb()
	cache.InitRedis()
	key := "video:comment_list:4396317246627840"

	emptyMember := "empty_member"
	z := redis.Z{
		Member: emptyMember,
		Score:  0,
	}
	count, err := cache.Redis.ZAdd(cache.Ctx, key, &z).Result()
	if err != nil {
		log.Print(err)
	}
	log.Print("创建数量", count)

	// 立即删除空成员
	count, err = cache.Redis.ZRem(cache.Ctx, key, emptyMember).Result()
	if err != nil {
		log.Print(err)
	}
	log.Print("删除数量", count)

	// 检查ZSET是否为空
	length, err := cache.Redis.ZCard(cache.Ctx, key).Result()
	if err != nil {
		log.Print(err)
	}

	if length == 0 {
		// ZSET为空
		println("ZSET为空")
	} else {
		// ZSET不为空
		println("ZSET不为空")
	}
}
