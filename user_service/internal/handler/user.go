package handler

import (
	"context"
	"user/internal/model"
	"user/internal/service"
)

type UserService struct {
}

func (*UserService) UserRegister(ctx context.Context, req *service.UserRequest) (resp *service.UserResponse, err error) {
	resp = new(service.UserResponse)
	var user model.User
	// Todo 在api层验证数据的有效性
	user.UserName = req.Username
	user.PassWord = req.Password

	err = model.GetInstance().Create(&user)

	return resp, nil
}
