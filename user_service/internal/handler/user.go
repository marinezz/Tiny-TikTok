package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user/internal/model"
	"user/internal/service"
	"utils/exception"
)

type UserService struct {
	service.UnimplementedUserServiceServer // 版本兼容问题
}

func NewUserService() *UserService {
	return &UserService{}
}

func (*UserService) UserRegister(ctx context.Context, req *service.UserRequest) (resp *service.UserResponse, err error) {
	resp = new(service.UserResponse)
	var user model.User

	// 检查用户是否已经存在
	if exist := model.GetInstance().CheckUserExist(req.Username); !exist {
		resp.StatusCode = exception.UserExist
		resp.StatusMsg = exception.GetMsg(exception.UserExist)
		resp.UserId = -1
		return resp, nil
	}
	// Todo 在api层验证数据的有效性
	user.UserName = req.Username
	user.PassWord = req.Password

	// 创建用户
	err = model.GetInstance().Create(&user)

	// 查询出ID
	userName, _ := model.GetInstance().FindUserByName(user.UserName)

	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	resp.UserId = userName.Id

	return resp, nil
}

func (*UserService) UserLogin(ctx context.Context, req *service.UserRequest) (*service.UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserLogin not implemented")
}
