package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"utils/exception"
	"video/internal/model"
	"video/internal/service"
	"video/pkg/cache"
)

// FavoriteAction 点赞操作
func (*VideoService) FavoriteAction(ctx context.Context, req *service.FavoriteActionRequest) (resp *service.FavoriteActionResponse, err error) {
	resp = new(service.FavoriteActionResponse)
	key := fmt.Sprintf("%s:%s", "user", "favorite_count")
	setKey := fmt.Sprintf("%s:%s:%s", "user", "favorite_video", strconv.FormatInt(req.UserId, 10))

	action := req.ActionType
	var favorite model.Favorite
	favorite.UserId = req.UserId
	favorite.VideoId = req.VideoId
	// 点赞操作
	if action == 1 {
		// 操作favorite表
		isAdd, err := model.GetFavoriteInstance().AddFavorite(&favorite)
		if err != nil {
			resp.StatusCode = exception.FavoriteErr
			resp.StatusMsg = exception.GetMsg(exception.FavoriteErr)
			return resp, err
		}

		// 操作video表，喜欢记录 + 1
		if isAdd == true {
			err := model.GetVideoInstance().AddFavoriteCount(req.VideoId)
			if err != nil {
				resp.StatusCode = exception.VideoFavoriteErr
				resp.StatusMsg = exception.GetMsg(exception.VideoFavoriteErr)
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
					return nil, fmt.Errorf("缓存错误：%v", err)
				}
			}

			// 加入喜欢set中，如果没有，构建缓存再加入set中
			exists, err := cache.Redis.Exists(cache.Ctx, setKey).Result()
			if err != nil {
				return nil, fmt.Errorf("缓存错误：%v", err)
			}

			if exists > 0 {
				err = cache.Redis.SAdd(cache.Ctx, setKey, strconv.FormatInt(req.VideoId, 10)).Err()
				if err != nil {
					return nil, fmt.Errorf("缓存错误：%v", err)
				}
			} else {
				err := buildFavoriteCache(req.UserId)
				if err != nil {
					return nil, fmt.Errorf("缓存错误：%v", err)
				}
				err = cache.Redis.SAdd(cache.Ctx, setKey, req.VideoId).Err()
				if err != nil {
					return nil, fmt.Errorf("缓存错误：%v", err)
				}
			}
		}
	}

	// 取消赞操作
	if action == 2 {
		// 操作favorite表
		err, isDelete := model.GetFavoriteInstance().DeleteFavorite(&favorite)
		if err != nil {
			resp.StatusCode = exception.CancelFavoriteErr
			resp.StatusMsg = exception.GetMsg(exception.CancelFavoriteErr)
			return resp, err
		}
		// 操作video表
		if isDelete == true {
			err := model.GetVideoInstance().DeleteFavoriteCount(req.VideoId)
			if err != nil {
				resp.StatusCode = exception.VideoFavoriteErr
				resp.StatusMsg = exception.GetMsg(exception.VideoFavoriteErr)
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
					return nil, fmt.Errorf("缓存错误：%v", err)
				}
			}

			// 加入喜欢set中，如果没有，构建缓存再去掉set中数据
			exists, err := cache.Redis.Exists(cache.Ctx, setKey).Result()
			if err != nil {
				return nil, fmt.Errorf("缓存错误：%v", err)
			}

			if exists > 0 {
				err = cache.Redis.SRem(cache.Ctx, setKey, req.VideoId).Err()
				if err != nil {
					return nil, fmt.Errorf("缓存错误：%v", err)
				}
			} else {
				err := buildFavoriteCache(req.UserId)
				if err != nil {
					return nil, fmt.Errorf("缓存错误：%v", err)
				}
				err = cache.Redis.SRem(context.Background(), setKey, req.VideoId).Err()
				if err != nil {
					return nil, fmt.Errorf("缓存错误：%v", err)
				}
			}
		}
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

// 构建点赞视频缓存
func buildFavoriteCache(userId int64) error {
	key := fmt.Sprintf("%s:%s:%s", "user", "favorite_video", strconv.FormatInt(userId, 10))

	// 查询出所有喜欢的视频
	favoriteVideoList, err := model.GetFavoriteInstance().FavoriteVideoList(userId)
	if err != nil {
		return err
	}

	videoIds := make([]interface{}, len(favoriteVideoList))
	for i, video := range favoriteVideoList {
		videoIds[i] = video
	}

	err = cache.Redis.SAdd(cache.Ctx, key, videoIds...).Err()
	if err != nil {
		return err
	}

	return nil
}
