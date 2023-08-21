package handler

import (
	"context"
	"social/internal/model"
	"social/internal/service"
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
	follow := model.Follow{
		UserId:   req.UserId,
		ToUserId: req.ToUserId,
		IsFollow: req.ActionType,
	}

	err = model.GetFollowInstance().FollowAction(&follow)
	if err != nil {
		resp.StatusCode = exception.ERROR
		resp.StatusMsg = exception.GetMsg(exception.ERROR)
		return resp, nil
	}
	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	return resp, nil
}

func (*SocialService) IsFollow(ctx context.Context, req *service.IsFollowRequest) (resp *service.IsFollowResponse, err error) {
	resp = new(service.IsFollowResponse)
	res, err := model.GetFollowInstance().IsFollow(req.UserId, req.ToUserId)
	if err != nil {
		resp.IsFollow = false
		resp.StatusCode = exception.ERROR
		resp.StatusMsg = exception.GetMsg(exception.ERROR)
		return resp, nil
	}
	resp.IsFollow = res
	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	return resp, nil
}

func (*SocialService) GetFollowList(ctx context.Context, req *service.FollowListRequest) (resp *service.FollowListResponse, err error) {
	resp = new(service.FollowListResponse)
	err = model.GetFollowInstance().GetFollowList(req.UserId, &resp.UserId)
	if err != nil {
		resp.StatusCode = exception.ERROR
		resp.StatusMsg = exception.GetMsg(exception.ERROR)
		return resp, nil
	}
	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	return resp, nil
}

func (*SocialService) GetFollowerList(ctx context.Context, req *service.FollowListRequest) (resp *service.FollowListResponse, err error) {
	resp = new(service.FollowListResponse)
	err = model.GetFollowInstance().GetFollowerList(req.UserId, &resp.UserId)
	if err != nil {
		resp.StatusCode = exception.ERROR
		resp.StatusMsg = exception.GetMsg(exception.ERROR)
		return resp, nil
	}
	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	return resp, nil
}

func (*SocialService) GetFriendList(ctx context.Context, req *service.FollowListRequest) (resp *service.FollowListResponse, err error) {
	resp = new(service.FollowListResponse)
	err = model.GetFollowInstance().GetFollowerList(req.UserId, &resp.UserId)
	if err != nil {
		resp.StatusCode = exception.ERROR
		resp.StatusMsg = exception.GetMsg(exception.ERROR)
		return resp, nil
	}
	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	return resp, nil
}

func (*SocialService) GetFollowCount(ctx context.Context, req *service.FollowCountRequest) (resp *service.FollowCountResponse, err error) {
	resp = new(service.FollowCountResponse)
	cnt, err := model.GetFollowInstance().GetFollowCount(req.UserId)
	if err != nil {
		resp.StatusCode = exception.ERROR
		resp.StatusMsg = exception.GetMsg(exception.ERROR)
		return resp, nil
	}
	resp.FollowCount = cnt
	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	return resp, nil
}

func (*SocialService) GetFollowerCount(ctx context.Context, req *service.FollowCountRequest) (resp *service.FollowerCountResponse, err error) {
	resp = new(service.FollowerCountResponse)
	cnt, err := model.GetFollowInstance().GetFollowerCount(req.UserId)
	if err != nil {
		resp.StatusCode = exception.ERROR
		resp.StatusMsg = exception.GetMsg(exception.ERROR)
		return resp, nil
	}
	resp.FollowerCount = cnt
	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	return resp, nil
}

func (*SocialService) GetFollowInfo(ctx context.Context, req *service.FollowInfoRequest) (resp *service.FollowInfoResponse, err error) {
	resp = new(service.FollowInfoResponse)
	for _, toUserId := range req.ToUserId {
		res1, err1 := model.GetFollowInstance().IsFollow(req.UserId, toUserId)
		cnt2, err2 := model.GetFollowInstance().GetFollowCount(toUserId)
		cnt3, err3 := model.GetFollowInstance().GetFollowerCount(toUserId)
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
	err = model.GetMessageInstance().GetMessage(req.UserId, req.ToUserId, &messages)
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
