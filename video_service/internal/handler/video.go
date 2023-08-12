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
