package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
	"utils/exception"
	"video/internal/model"
	"video/internal/service"
	"video/pkg/cache"
)

// FavoriteAction 点赞操作 todo 也可以设计成定时任务
func (*VideoService) FavoriteAction(ctx context.Context, req *service.FavoriteActionRequest) (resp *service.FavoriteActionResponse, err error) {
	resp = new(service.FavoriteActionResponse)
	key := fmt.Sprintf("%s:%s", "user", "favorite_count")
	setKey := fmt.Sprintf("%s:%s:%s", "video", "favorite_video", strconv.FormatInt(req.VideoId, 10))
	favoriteKey := fmt.Sprintf("%s:%s:%s", "user", "favorit_list", strconv.FormatInt(req.UserId, 10))

	action := req.ActionType
	var favorite model.Favorite
	favorite.UserId = req.UserId
	favorite.VideoId = req.VideoId

	// 查看缓存是否存在，不存在这构建一次缓存，避免极端情况
	setExists, err := cache.Redis.Exists(cache.Ctx, setKey).Result()
	if err != nil {
		return nil, fmt.Errorf("缓存错误：%v", err)
	}
	if setExists == 0 {
		err := buildVideoFavorite(req.VideoId)
		if err != nil {
			return nil, fmt.Errorf("缓存错误：%v", err)
		}
	}

	// 点赞操作
	if action == 1 {
		// 查看缓存，避免重复点赞
		result, err := cache.Redis.SIsMember(cache.Ctx, setKey, req.UserId).Result()
		if err != nil {
			return nil, fmt.Errorf("缓存错误：%v", err)
		}

		if result {
			// 重复点赞
			resp.StatusCode = exception.FavoriteErr
			resp.StatusMsg = exception.GetMsg(exception.FavoriteErr)
			return resp, err
		}

		// 操作favorite表
		tx := model.DB.Begin()
		err = model.GetFavoriteInstance().AddFavorite(tx, &favorite)
		if err != nil {
			//tx.Rollback()
			resp.StatusCode = exception.FavoriteErr
			resp.StatusMsg = exception.GetMsg(exception.FavoriteErr)
			return resp, err
		}

		// 点赞成功，缓存中点赞总数 + 1
		exist, err := cache.Redis.HExists(cache.Ctx, key, strconv.FormatInt(req.UserId, 10)).Result()
		if err != nil {
			return nil, fmt.Errorf("缓存错误：%v", err)
		}

		if exist {
			// 字段存在，该记录数量 + 1
			_, err = cache.Redis.HIncrBy(cache.Ctx, key, strconv.FormatInt(req.UserId, 10), 1).Result()
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("缓存错误：%v", err)
			}
		}

		// 加入喜欢set中，如果没有，构建缓存再加入set中
		err = cache.Redis.SAdd(cache.Ctx, setKey, strconv.FormatInt(req.UserId, 10)).Err()
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("缓存错误：%v", err)
		}

		// 删除喜欢列表缓存
		err = cache.Redis.Del(cache.Ctx, favoriteKey).Err()
		if err != nil {
			return nil, fmt.Errorf("缓存错误：%v", err)
		}
		defer func() {
			go func() {
				//延时3秒执行
				time.Sleep(time.Second * 3)
				//再次删除缓存
				cache.Redis.Del(cache.Ctx, favoriteKey)
			}()
		}()

		tx.Commit()
	}

	// 取消赞操作
	if action == 2 {
		// 查看缓存，避免重复点删除
		result, err := cache.Redis.SIsMember(cache.Ctx, setKey, req.UserId).Result()
		if err != nil {
			return nil, fmt.Errorf("缓存错误：%v", err)
		}

		if result == false {
			// 重复删除
			resp.StatusCode = exception.CancelFavoriteErr
			resp.StatusMsg = exception.GetMsg(exception.CancelFavoriteErr)
			return resp, err
		}

		// 操作favorite表
		tx := model.DB.Begin()
		err = model.GetFavoriteInstance().DeleteFavorite(tx, &favorite)
		if err != nil {
			resp.StatusCode = exception.CancelFavoriteErr
			resp.StatusMsg = exception.GetMsg(exception.CancelFavoriteErr)
			return resp, err
		}

		// 点赞成功，缓存中点赞总数 - 1
		exist, err := cache.Redis.HExists(cache.Ctx, key, strconv.FormatInt(req.UserId, 10)).Result()
		if err != nil {
			return nil, fmt.Errorf("缓存错误：%v", err)
		}

		if exist {
			// 字段存在，该记录数量 + 1
			_, err = cache.Redis.HIncrBy(cache.Ctx, key, strconv.FormatInt(req.UserId, 10), -1).Result()
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("缓存错误：%v", err)
			}
		}

		// set中删除
		err = cache.Redis.SRem(cache.Ctx, setKey, req.UserId).Err()
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("缓存错误：%v", err)
		}

		// 删除喜欢列表缓存
		err = cache.Redis.Del(cache.Ctx, favoriteKey).Err()
		if err != nil {
			return nil, fmt.Errorf("缓存错误：%v", err)
		}
		defer func() {
			go func() {
				//延时3秒执行
				time.Sleep(time.Second * 3)
				//再次删除缓存
				cache.Redis.Del(cache.Ctx, favoriteKey)
			}()
		}()

		tx.Commit()
	}

	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)

	return resp, nil
}

// FavoriteList 喜欢列表
func (*VideoService) FavoriteList(ctx context.Context, req *service.FavoriteListRequest) (resp *service.FavoriteListResponse, err error) {
	resp = new(service.FavoriteListResponse)
	var videos []model.Video
	key := fmt.Sprintf("%s:%s:%s", "user", "favorit_list", strconv.FormatInt(req.UserId, 10))

	exits, err := cache.Redis.Exists(cache.Ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("缓存错误：%v", err)
	}

	if exits > 0 {
		videosString, err := cache.Redis.Get(cache.Ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("缓存错误：%v", err)
		}
		err = json.Unmarshal([]byte(videosString), &videos)
		if err != nil {
			return nil, err
		}
	} else {
		// 根据用户id找到所有的视频
		var videoIds []int64
		videoIds, err = model.GetFavoriteInstance().FavoriteVideoList(req.UserId)
		if err != nil {
			resp.StatusCode = exception.UserNoVideo
			resp.StatusMsg = exception.GetMsg(exception.UserNoVideo)
			return resp, err
		}

		// 根据视频id找到视频的详细信息
		videos, err = model.GetVideoInstance().GetVideoList(videoIds)
		if err != nil {
			resp.StatusCode = exception.VideoUnExist
			resp.StatusMsg = exception.GetMsg(exception.VideoUnExist)
			return resp, err
		}

		// 放入缓存中
		videosJson, _ := json.Marshal(videos)
		err := cache.Redis.Set(cache.Ctx, key, videosJson, 30*time.Minute).Err()
		if err != nil {
			return nil, fmt.Errorf("缓存错误：%v", err)
		}
	}

	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	resp.VideoList = BuildVideoForFavorite(videos, true)

	return resp, nil
}

// 查询缓存，判断是否喜欢
func isFavorite(userId int64, videoId int64) bool {
	var isFavorite bool
	key := fmt.Sprintf("%s:%s:%s", "video", "favorite_video", strconv.FormatInt(videoId, 10))

	exists, err := cache.Redis.Exists(cache.Ctx, key).Result()
	if err != nil {
		log.Print(err)
	}

	if exists > 0 {
		isFavorite, err = cache.Redis.SIsMember(cache.Ctx, key, strconv.FormatInt(userId, 10)).Result()
		if err != nil {
			log.Print(err)
		}
	} else {
		err := buildVideoFavorite(videoId)
		if err != nil {
			log.Print(err)
		}
		isFavorite, err = cache.Redis.SIsMember(cache.Ctx, key, strconv.FormatInt(userId, 10)).Result()
		if err != nil {
			log.Print(err)
		}
	}

	return isFavorite
}

// 构建视频点赞缓存
func buildVideoFavorite(videoId int64) error {
	key := fmt.Sprintf("%s:%s:%s", "video", "favorite_video", strconv.FormatInt(videoId, 10))

	// 查询出所有喜欢的视频
	userIdList, err := model.GetFavoriteInstance().FavoriteUserList(videoId)
	if err != nil {
		return err
	}

	// 如果点赞数量为空，则不会创建cache，所以设计一个先放入，再删除，创建一个空记录。避免反复查表
	userIds := make([]interface{}, len(userIdList))
	for i, video := range userIdList {
		userIds[i] = video
	}

	err = cache.Redis.SAdd(cache.Ctx, key, userIds...).Err()
	if err != nil {
		return err
	}

	return nil
}

// 通过缓存查询视频的获赞数量
func getFavoriteCount(videoId int64) int64 {
	setKey := fmt.Sprintf("%s:%s:%s", "video", "favorite_video", strconv.FormatInt(videoId, 10))

	// 查看缓存是否存在，不存在这构建一次缓存，避免极端情况
	setExists, err := cache.Redis.Exists(cache.Ctx, setKey).Result()
	if err != nil {
		log.Print(err)
	}

	if setExists == 0 {
		err := buildVideoFavorite(videoId)
		if err != nil {
			log.Print(err)
		}
	}

	count, err := cache.Redis.SCard(cache.Ctx, setKey).Result()
	if err != nil {
		log.Print(err)
	}
	return count
}
