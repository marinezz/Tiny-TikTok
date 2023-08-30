package handler

import (
	"context"
	"strconv"
	"utils/exception"
	"video/internal/model"
	"video/internal/service"
	"video/pkg/cache"
)

// FavoriteAction 点赞操作
func (*VideoService) FavoriteAction(ctx context.Context, req *service.FavoriteActionRequest) (resp *service.FavoriteActionResponse, err error) {
	resp = new(service.FavoriteActionResponse)

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
			exist, err := cache.Redis.HExists(context.Background(), "userFavorite_count", strconv.FormatInt(req.UserId, 10)).Result()
			if err != nil {
				resp.StatusCode = exception.CacheErr
				resp.StatusMsg = exception.GetMsg(exception.CacheErr)
				return resp, err
			}

			if exist {
				// 字段存在，该记录数量 + 1
				_, err = cache.Redis.HIncrBy(context.Background(), "userFavorite_count", strconv.FormatInt(req.UserId, 10), 1).Result()
				if err != nil {
					resp.StatusCode = exception.CacheErr
					resp.StatusMsg = exception.GetMsg(exception.CacheErr)
					return resp, err
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
		}

		// 点赞成功，缓存中点赞总数 + 1
		exist, err := cache.Redis.HExists(context.Background(), "userFavorite_count", strconv.FormatInt(req.UserId, 10)).Result()
		if err != nil {
			resp.StatusCode = exception.CacheErr
			resp.StatusMsg = exception.GetMsg(exception.CacheErr)
			return resp, err
		}

		if exist {
			// 字段存在，该记录数量 + 1
			_, err = cache.Redis.HIncrBy(context.Background(), "userFavorite_count", strconv.FormatInt(req.UserId, 10), -1).Result()
			if err != nil {
				resp.StatusCode = exception.CacheErr
				resp.StatusMsg = exception.GetMsg(exception.CacheErr)
				return resp, err
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

	// 根据用户id找到所有的视频
	var videoIds []int64
	videoIds, err = model.GetFavoriteInstance().FavoriteVideoList(req.UserId)
	if err != nil {
		resp.StatusCode = exception.UserNoVideo
		resp.StatusMsg = exception.GetMsg(exception.UserNoVideo)
		return resp, err
	}

	// 根据视频id找到视频的详细信息
	videos, err := model.GetVideoInstance().GetVideoList(videoIds)
	if err != nil {
		resp.StatusCode = exception.VideoUnExist
		resp.StatusMsg = exception.GetMsg(exception.VideoUnExist)
		return resp, err
	}

	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	resp.VideoList = BuildVideoForFavorite(videos, true)

	return resp, nil
}
