package handler

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*VideoService) Feed(ctx context.Context, req *service.FeedRequest) (resp *service.FeedResponse, err error) {

	return nil, status.Errorf(codes.Unimplemented, "method Feed not implemented")
}

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

func (*VideoService) CommentAction(ctx context.Context, req *service.CommentActionRequest) (resp *service.CommentActionResponse, err error) {
	resp = new(service.CommentActionResponse)
	comment := model.Comment{
		UserId:  req.UserId,
		VideoId: req.VideoId,
		Content: req.CommentText,
	}
	action := req.ActionType

	time := time.Now()

	// 发布评论
	if action == 1 {
		comment.CreatAt = time
		id, _ := model.GetCommentInstance().CreateComment(&comment)

		// 视频评论数量 + 1
		model.GetVideoInstance().AddCommentCount(req.VideoId)

		commentResp := &service.Comment{
			Id:      id,
			Content: req.CommentText,
			// 将Time.time转换成字符串形式
			CreateDate: time.Format("2006-01-02 15:04:05"),
		}

		// 将评论返回
		resp.StatusCode = exception.SUCCESS
		resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
		resp.Comment = commentResp

		return resp, nil
	}

	// 删除评论
	model.GetCommentInstance().DeleteComment(req.CommentId)
	// 视频评论数量 - 1
	model.GetVideoInstance().DeleteCommentCount(req.VideoId)

	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	resp.Comment = nil

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
