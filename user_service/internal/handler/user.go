package handler

import (
	"context"
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

func (*UserService) UserLogin(ctx context.Context, req *service.UserRequest) (resp *service.UserResponse, err error) {
	resp = new(service.UserResponse)

	// 检查用户是否存在
	if exist := model.GetInstance().CheckUserExist(req.Username); exist {
		resp.StatusCode = exception.UserUnExist
		resp.StatusMsg = exception.GetMsg(exception.UserUnExist)
		resp.UserId = -1
		return resp, nil
	}

	// 检查密码是否正确
	user, err := model.GetInstance().FindUserByName(req.Username)
	if ok := model.GetInstance().CheckPassWord(req.Password, user.PassWord); !ok {
		resp.StatusCode = exception.PasswordError
		resp.StatusMsg = exception.GetMsg(exception.PasswordError)
		resp.UserId = -1
		return resp, nil
	}

	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	resp.UserId = user.Id

	return resp, nil
}
