package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"user/internal/model"
	"user/internal/service"
	"user/pkg/cache"
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
		return resp, err
	}

	user.UserName = req.Username
	user.PassWord = req.Password

	// 创建用户
	err = model.GetInstance().Create(&user)
	if err != nil {
		resp.StatusCode = exception.DataErr
		resp.StatusMsg = exception.GetMsg(exception.DataErr)
		return resp, err
	}

	// 查询出ID
	userName, err := model.GetInstance().FindUserByName(user.UserName)
	if err != nil {
		resp.StatusCode = exception.UserUnExist
		resp.StatusMsg = exception.GetMsg(exception.UserUnExist)
		return resp, err
	}

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
		return resp, err
	}

	// 检查密码是否正确
	user, err := model.GetInstance().FindUserByName(req.Username)
	if ok := model.GetInstance().CheckPassWord(req.Password, user.PassWord); !ok {
		resp.StatusCode = exception.PasswordError
		resp.StatusMsg = exception.GetMsg(exception.PasswordError)
		return resp, err
	}

	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)
	resp.UserId = user.Id

	return resp, nil
}

// UserInfo 用户信息
func (*UserService) UserInfo(ctx context.Context, req *service.UserInfoRequest) (resp *service.UserInfoResponse, err error) {
	resp = new(service.UserInfoResponse)

	// 根据userId切片查询user信息
	userIds := req.UserIds

	for _, userId := range userIds {
		// 查看缓存是否存在 需要的信息
		var user *model.User
		key := fmt.Sprintf("%s:%s:%s", "user", "info", strconv.FormatInt(userId, 10))

		count, err := cache.Redis.Exists(context.Background(), key).Result()
		if err != nil {
			resp.StatusCode = exception.CacheErr
			resp.StatusMsg = exception.GetMsg(exception.CacheErr)
			return resp, err
		}

		if count > 0 {
			// 查询缓存
			userString, err := cache.Redis.Get(context.Background(), key).Result()
			if err != nil {
				resp.StatusCode = exception.CacheErr
				resp.StatusMsg = exception.GetMsg(exception.CacheErr)
				return resp, err
			}
			err = json.Unmarshal([]byte(userString), &user)
			if err != nil {
				return nil, err
			}
		} else {
			// 查询数据库
			user, err = model.GetInstance().FindUserById(userId)
			if err != nil {
				resp.StatusCode = exception.UserUnExist
				resp.StatusMsg = exception.GetMsg(exception.UserUnExist)
				return resp, err
			}

			// 将查询结果放入缓存中
			userJson, _ := json.Marshal(&user)
			err = cache.Redis.Set(context.Background(), key, userJson, 12*time.Hour).Err()
			if err != nil {
				resp.StatusCode = exception.CacheErr
				resp.StatusMsg = exception.GetMsg(exception.CacheErr)
				return resp, err
			}

		}
		resp.Users = append(resp.Users, BuildUser(user))
	}

	resp.StatusCode = exception.SUCCESS
	resp.StatusMsg = exception.GetMsg(exception.SUCCESS)

	return resp, nil
}

func BuildUser(u *model.User) *service.User {
	user := service.User{
		Id:              u.Id,
		Name:            u.UserName,
		Avatar:          u.Avatar,
		BackgroundImage: u.BackgroundImage,
		Signature:       u.Signature,
	}
	return &user
}
