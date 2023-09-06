package handler

import (
	"context"
	"social/internal/model"
	"social/internal/service"
	"social/pkg/cache"
	"utils/exception"
)

type SocialService struct {
	service.UnimplementedSocialServiceServer // 版本兼容问题
}

func NewSocialService() *SocialService {
	return &SocialService{}
}

// FollowAction 关注服务
func (*SocialService) FollowAction(ctx context.Context, req *service.FollowRequest) (resp *service.FollowResponse, err error) {
	resp = new(service.FollowResponse)

	if req.UserId == req.ToUserId {
		resp.StatusCode = exception.FollowSelfErr
		resp.StatusMsg = exception.GetMsg(exception.FollowSelfErr)
		return resp, nil
	}

	resp.StatusCode = exception.ERROR
	resp.StatusMsg = exception.GetMsg(exception.ERROR)

	/* mysql操作

	follow := model.Follow{
		UserId:   req.UserId,
		ToUserId: req.ToUserId,
		IsFollow: req.ActionType,
	}

	err = model.GetFollowInstance().FollowAction(&follow)
	*/

	// redis操作
	err = cache.FollowAction(req.UserId, req.ToUserId, req.ActionType)

	if err != nil {
		return resp, nil
	}
	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	return resp, nil
}

func (*SocialService) GetFollowList(ctx context.Context, req *service.FollowListRequest) (resp *service.FollowListResponse, err error) {
	resp = new(service.FollowListResponse)
	resp.StatusCode = exception.ERROR
	resp.StatusMsg = exception.GetMsg(exception.ERROR)

	// mysql操作
	// err = model.GetFollowInstance().GetFollowList(req.UserId, &resp.UserId)

	// redis操作
	err = cache.GetFollowList(req.UserId, &resp.UserId)

	if err != nil {
		return resp, nil
	}
	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	return resp, nil
}

func (*SocialService) GetFollowerList(ctx context.Context, req *service.FollowListRequest) (resp *service.FollowListResponse, err error) {
	resp = new(service.FollowListResponse)
	resp.StatusCode = exception.ERROR
	resp.StatusMsg = exception.GetMsg(exception.ERROR)

	// mysql操作
	// err = model.GetFollowInstance().GetFollowerList(req.UserId, &resp.UserId)

	// redis操作
	err = cache.GetFollowerList(req.UserId, &resp.UserId)
	if err != nil {
		return resp, nil
	}
	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	return resp, nil
}

func (*SocialService) GetFriendList(ctx context.Context, req *service.FollowListRequest) (resp *service.FollowListResponse, err error) {
	resp = new(service.FollowListResponse)
	resp.StatusCode = exception.ERROR
	resp.StatusMsg = exception.GetMsg(exception.ERROR)

	// mysql
	// err = model.GetFollowInstance().GetFriendList(req.UserId, &resp.UserId)

	// mysql
	err = cache.GetFriendList(req.UserId, &resp.UserId)

	if err != nil {
		return resp, nil
	}
	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	return resp, nil
}

func (*SocialService) GetFollowInfo(ctx context.Context, req *service.FollowInfoRequest) (resp *service.FollowInfoResponse, err error) {
	resp = new(service.FollowInfoResponse)
	for _, toUserId := range req.ToUserId {
		/* mysql
		res1, err1 := model.GetFollowInstance().IsFollow(req.UserId, toUserId)
		cnt2, err2 := model.GetFollowInstance().GetFollowCount(toUserId)
		cnt3, err3 := model.GetFollowInstance().GetFollowerCount(toUserId)
		*/
		res1, err1 := cache.IsFollow(req.UserId, toUserId)
		cnt2, err2 := cache.GetFollowCount(toUserId)
		cnt3, err3 := cache.GetFollowerCount(toUserId)
		if err1 != nil || err2 != nil || err3 != nil {
			resp.StatusCode = exception.ERROR
			resp.StatusMsg = exception.GetMsg(exception.ERROR)
			return resp, nil
		}
		resp.FollowInfo = append(resp.FollowInfo, &service.FollowInfo{
			IsFollow:      res1,
			FollowCount:   cnt2,
			FollowerCount: cnt3,
			ToUserId:      toUserId,
		})
	}
	return resp, nil
}

// PostMessage 消息服务
func (*SocialService) PostMessage(ctx context.Context, req *service.PostMessageRequest) (resp *service.PostMessageResponse, err error) {
	resp = new(service.PostMessageResponse)
	message := model.Message{
		UserId:   req.UserId,
		ToUserId: req.ToUserId,
		Message:  req.Content,
	}
	err = model.GetMessageInstance().PostMessage(&message)
	if err != nil {
		resp.StatusCode = exception.ERROR
		resp.StatusMsg = exception.GetMsg(exception.ERROR)
		return resp, nil
	}
	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	return resp, nil
}

func (*SocialService) GetMessage(ctx context.Context, req *service.GetMessageRequest) (resp *service.GetMessageResponse, err error) {
	resp = new(service.GetMessageResponse)
	var messages []model.Message
	err = model.GetMessageInstance().GetMessage(req.UserId, req.ToUserId, req.PreMsgTime, &messages)
	if err != nil {
		resp.StatusCode = exception.ERROR
		resp.StatusMsg = exception.GetMsg(exception.ERROR)
		return resp, nil
	}

	for _, message := range messages {
		resp.Message = append(resp.Message, &service.Message{
			Id:        message.Id,
			UserId:    message.UserId,
			ToUserId:  message.ToUserId,
			Content:   message.Message,
			CreatedAt: message.CreatedAt,
		})
	}
	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	return resp, nil
}
