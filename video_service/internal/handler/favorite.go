package handler

import (
	"context"
	"utils/exception"
	"video/internal/model"
	"video/internal/service"
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
		model.GetFavoriteInstance().AddFavorite(&favorite)
		// 操作video表，喜欢记录 + 1
		model.GetVideoInstance().AddFavoriteCount(req.VideoId)
	}

	// 取消赞操作
	if action == 2 {
		// 操作favorite表
		model.GetFavoriteInstance().DeleteFavorite(&favorite)
		// 操作video表
		model.GetVideoInstance().DeleteFavoriteCount(req.VideoId)
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
	videoIds, _ = model.GetFavoriteInstance().FavoriteVideoList(req.UserId)

	// 根据视频id找到视频的详细信息
	videos, _ := model.GetVideoInstance().GetVideoList(videoIds)

	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	resp.VideoList = BuildVideoForFavorite(videos, true)

	return resp, nil
}
