package handler

import (
	"context"
	"github.com/google/uuid"
	"log"
	"time"
	"utils/exception"
	"video/internal/model"
	"video/internal/service"
	"video/pkg/cut"
	"video/third_party"
)

type VideoService struct {
	service.UnimplementedVideoServiceServer // 版本兼容问题
}

func NewVideoService() *VideoService {
	return &VideoService{}
}

// Feed 视频流
func (*VideoService) Feed(ctx context.Context, req *service.FeedRequest) (resp *service.FeedResponse, err error) {
	resp = new(service.FeedResponse)

	// 获取时间
	var timePoint time.Time
	if req.LatestTime == -1 {
		timePoint = time.Now()
	} else {
		timePoint = time.Unix(req.LatestTime, 0)
	}

	// 根据时间获取视频
	videos, _ := model.GetVideoInstance().GetVideoByTime(timePoint)

	if req.UserId == -1 {
		// 用户没有登录
		resp.VideoList = BuildVideoForFavorite(videos, false)
	} else {
		resp.VideoList = BuildVideo(videos, req.UserId)
	}

	// 获取列表中最早发布视频的时间作为下一次请求的时间
	LastIndex := len(videos) - 1
	resp.NextTime = videos[LastIndex].CreatAt.Unix()

	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)

	return resp, nil
}

// PublishAction 发布视频
func (*VideoService) PublishAction(ctx context.Context, req *service.PublishActionRequest) (resp *service.PublishActionResponse, err error) {
	resp = new(service.PublishActionResponse)

	// 获取参数,生成地址
	title := req.Title
	UUID := uuid.New()
	videoDir := title + "--" + UUID.String() + ".mp4"
	pictureDir := title + "--" + UUID.String() + ".jpg"

	videoUrl := "http://tiny-tiktok.oss-cn-chengdu.aliyuncs.com/" + videoDir
	pictureUrl := "http://tiny-tiktok.oss-cn-chengdu.aliyuncs.com/" + pictureDir

	// 上传视频，切取封面，上传图片
	go func() {
		// 上传视频
		third_party.Upload(videoDir, req.Data)
		// 获取封面,获取第几秒的封面
		coverByte, _ := cut.Cover(videoUrl, "00:00:01")
		// 上传封面
		third_party.Upload(pictureDir, coverByte)
		log.Print("上传成功")
	}()
	// 创建video
	video := model.Video{
		AuthId:        req.UserId,
		Title:         title,
		CoverUrl:      pictureUrl,
		PlayUrl:       videoUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		CreatAt:       time.Now(),
	}
	err = model.GetVideoInstance().Create(&video)
	if err != nil {
		resp.StatusCode = exception.ERROR
		resp.StatusMsg = exception.GetMsg(exception.ERROR)
		return resp, nil
	}
	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)

	return resp, nil
}

// PublishList 发布列表
func (*VideoService) PublishList(ctx context.Context, req *service.PublishListRequest) (resp *service.PublishListResponse, err error) {
	resp = new(service.PublishListResponse)

	// 根据用户id找到所有的视频
	videos, _ := model.GetVideoInstance().GetVideoListByUser(req.UserId)

	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	resp.VideoList = BuildVideo(videos, req.UserId)

	return resp, nil
}

// CountInfo 计数信息
func (*VideoService) CountInfo(ctx context.Context, req *service.CountRequest) (resp *service.CountResponse, err error) {
	resp = new(service.CountResponse)

	userIds := req.UserIds

	for _, userId := range userIds {
		var count service.Count
		// 获取赞的数量
		count.TotalFavorited, _ = model.GetVideoInstance().GetFavoritedCount(userId)
		// 获取作品数量
		count.WorkCount, _ = model.GetVideoInstance().GetWorkCount(userId)
		// 获取喜欢数量
		count.FavoriteCount, _ = model.GetFavoriteInstance().GetFavoriteCount(userId)

		resp.Counts = append(resp.Counts, &count)
	}

	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)

	return resp, nil
}

func BuildVideo(videos []model.Video, userId int64) []*service.Video {
	var videoResp []*service.Video

	for _, video := range videos {
		// 获取is_favorite的状态
		isFavorite, _ := model.GetFavoriteInstance().IsFavorite(userId, video.Id)

		videoResp = append(videoResp, &service.Video{
			Id:            video.Id,
			AuthId:        video.AuthId,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    isFavorite,
			Title:         video.Title,
		})
	}

	return videoResp
}

func BuildVideoForFavorite(videos []model.Video, isFavorite bool) []*service.Video {
	var videoResp []*service.Video

	for _, video := range videos {
		videoResp = append(videoResp, &service.Video{
			Id:            video.Id,
			AuthId:        video.AuthId,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    isFavorite,
			Title:         video.Title,
		})
	}

	return videoResp
}
