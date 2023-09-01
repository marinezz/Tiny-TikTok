package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"strconv"
	"sync"
	"time"
	"utils/exception"
	"video/internal/model"
	"video/internal/service"
	"video/pkg/cache"
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
		timePoint = time.Unix(req.LatestTime/1000, 0)
	}

	// 根据时间获取视频
	videos, err := model.GetVideoInstance().GetVideoByTime(timePoint)
	if err != nil {
		resp.StatusCode = exception.VideoUnExist
		resp.StatusMsg = exception.GetMsg(exception.VideoUnExist)
		return resp, err
	}

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
	var updataErr, creatErr error
	resp = new(service.PublishActionResponse)
	key := fmt.Sprintf("%s:%s", "user", "work_count")

	// 获取参数,生成地址
	title := req.Title
	UUID := uuid.New()
	videoDir := title + "--" + UUID.String() + ".mp4"
	pictureDir := title + "--" + UUID.String() + ".jpg"

	videoUrl := "http://tiny-tiktok.oss-cn-chengdu.aliyuncs.com/" + videoDir
	pictureUrl := "http://tiny-tiktok.oss-cn-chengdu.aliyuncs.com/" + pictureDir

	// 等待上传和创建数组库完成
	var wg sync.WaitGroup
	wg.Add(2)

	// 上传视频，切取封面，上传图片
	go func() {
		defer wg.Done()
		// 上传视频
		updataErr = third_party.Upload(videoDir, req.Data)
		// 获取封面,获取第几秒的封面
		coverByte, _ := cut.Cover(videoUrl, "00:00:03")
		// 上传封面
		updataErr = third_party.Upload(pictureDir, coverByte)
		log.Print("上传成功")
	}()

	// 创建数据
	go func() {
		defer wg.Done()
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
		creatErr = model.GetVideoInstance().Create(&video)
	}()

	wg.Wait()

	// 异步回滚
	if updataErr != nil || creatErr != nil {
		go func() {
			// 存入数据库失败，删除上传
			if creatErr != nil {
				_ = third_party.Delete(videoDir)
				_ = third_party.Delete(pictureDir)
			}
			// 上传失败，删除数据库
			if updataErr != nil {
				// TODO 根据url查找，效率比较低
				_ = model.GetVideoInstance().DeleteVideoByUrl(videoUrl)
			}
		}()
	}
	if updataErr != nil || creatErr != nil {
		resp.StatusCode = exception.VideoUploadErr
		resp.StatusMsg = exception.GetMsg(exception.VideoUploadErr)
		return resp, updataErr
	}

	// 发布成功，缓存中作品总数 + 1，如果不存在缓存则不做操作
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

	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)

	return resp, nil
}

// PublishList 发布列表
func (*VideoService) PublishList(ctx context.Context, req *service.PublishListRequest) (resp *service.PublishListResponse, err error) {
	resp = new(service.PublishListResponse)
	var videos []model.Video
	key := fmt.Sprintf("%s:%s:%s", "user", "work_list", strconv.FormatInt(req.UserId, 10))

	// 根据用户id找到所有的视频,先找缓存，再查数据库
	exists, err := cache.Redis.Exists(cache.Ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("缓存错误：%v", err)
	}

	if exists > 0 {
		videosString, err := cache.Redis.Get(cache.Ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("缓存错误：%v", err)
		}
		err = json.Unmarshal([]byte(videosString), &videos)
		if err != nil {
			return nil, err
		}
	} else {
		videos, err = model.GetVideoInstance().GetVideoListByUser(req.UserId)
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
		exist, err := cache.Redis.HExists(cache.Ctx, "user:total_favorite", strconv.FormatInt(userId, 10)).Result()
		if err != nil {
			return nil, fmt.Errorf("缓存错误：%v", err)
		}

		if exist {
			count.FavoriteCount, err = cache.Redis.HGet(cache.Ctx, "user:total_favorite", strconv.FormatInt(userId, 10)).Int64()
			if err != nil {
				return nil, fmt.Errorf("缓存错误：%v", err)
			}
		} else {
			count.TotalFavorited, err = model.GetVideoInstance().GetFavoritedCount(userId)
			if err != nil {
				resp.StatusCode = exception.VideoNoFavorite
				resp.StatusMsg = exception.GetMsg(exception.VideoNoFavorite)
				return resp, err
			}

			// 放入缓存
			err := cache.Redis.HSet(cache.Ctx, "user:total_favorite", strconv.FormatInt(userId, 10), count.TotalFavorited).Err()
			if err != nil {
				return nil, fmt.Errorf("缓存错误：%v", err)
			}
		}

		// 获取作品数量
		exist, err = cache.Redis.HExists(cache.Ctx, "user:work_count", strconv.FormatInt(userId, 10)).Result()
		if err != nil {
			return nil, fmt.Errorf("缓存错误：%v", err)
		}
		// 如果存在则读缓存
		if exist {
			count.WorkCount, err = cache.Redis.HGet(cache.Ctx, "user:work_count", strconv.FormatInt(userId, 10)).Int64()
			if err != nil {
				return nil, fmt.Errorf("缓存错误：%v", err)
			}
		} else {
			// 不存在则查数据库
			count.WorkCount, err = model.GetVideoInstance().GetWorkCount(userId)
			if err != nil {
				resp.StatusCode = exception.UserNoVideo
				resp.StatusMsg = exception.GetMsg(exception.UserNoVideo)
				return resp, err
			}
			// 放入缓存
			err := cache.Redis.HSet(cache.Ctx, "user:work_count", strconv.FormatInt(userId, 10), count.WorkCount).Err()
			if err != nil {
				return nil, fmt.Errorf("缓存错误：%v", err)
			}
		}

		// 获取喜欢数量
		exist, err = cache.Redis.HExists(cache.Ctx, "user:favorite_count", strconv.FormatInt(userId, 10)).Result()
		if err != nil {
			return nil, fmt.Errorf("缓存错误：%v", err)
		}
		if exist {
			count.WorkCount, err = cache.Redis.HGet(cache.Ctx, "user:favorite_count", strconv.FormatInt(userId, 10)).Int64()
			if err != nil {
				return nil, fmt.Errorf("缓存错误：%v", err)
			}
		} else {
			count.FavoriteCount, err = model.GetFavoriteInstance().GetFavoriteCount(userId)
			if err != nil {
				resp.StatusCode = exception.UserNoFavorite
				resp.StatusMsg = exception.GetMsg(exception.UserNoFavorite)
				return resp, err
			}

			// 放入缓存
			err := cache.Redis.HSet(cache.Ctx, "user:favorite_count", strconv.FormatInt(userId, 10), count.FavoriteCount).Err()
			if err != nil {
				return nil, fmt.Errorf("缓存错误：%v", err)
			}
		}

		resp.Counts = append(resp.Counts, &count)
	}

	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)

	return resp, nil
}

func BuildVideo(videos []model.Video, userId int64) []*service.Video {
	var videoResp []*service.Video

	for _, video := range videos {
		// 查询是否有喜欢的缓存，如果有，比对缓存，如果没有，构建缓存再查缓存
		var isFavorite bool
		key := fmt.Sprintf("%s:%s:%s", "user", "favorite_video", strconv.FormatInt(userId, 10))

		exists, err := cache.Redis.Exists(cache.Ctx, key).Result()
		if err != nil {
			log.Print(err)
		}

		if exists > 0 {
			isFavorite, err = cache.Redis.SIsMember(cache.Ctx, key, strconv.FormatInt(video.Id, 10)).Result()
			if err != nil {
				log.Print(err)
			}
		} else {
			err := buildFavoriteCache(userId)
			if err != nil {
				log.Print(err)
			}
			isFavorite, err = cache.Redis.SIsMember(cache.Ctx, key, strconv.FormatInt(video.Id, 10)).Result()
			if err != nil {
				log.Print(err)
			}
		}

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
