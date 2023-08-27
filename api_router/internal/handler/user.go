package handler

import (
	"api_router/internal/service"
	"api_router/pkg/auth"
	"api_router/pkg/res"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
	"utils/exception"
)

// UserRegister 用户注册
func UserRegister(ctx *gin.Context) {
	var userReq service.UserRequest

	userReq.Username = ctx.Query("username")
	userReq.Password = ctx.Query("password")

	userServiceClient := ctx.Keys["user_service"].(service.UserServiceClient)
	userResp, err := userServiceClient.UserRegister(context.Background(), &userReq)
	if err != nil {
		PanicIfUserError(err)
	}
	token, _ := auth.GenerateToken(userResp.UserId)

	r := res.UserResponse{
		StatusCode: userResp.StatusCode,
		StatusMsg:  userResp.StatusMsg,
		UserId:     userResp.UserId,
		Token:      token,
	}

	ctx.JSON(http.StatusOK, r)
}

// UserLogin 用户登录
func UserLogin(ctx *gin.Context) {
	var userReq service.UserRequest
	//err := ctx.Bind(&userReq)
	//if err != nil {
	//	PanicIfUserError(err)
	//}

	userReq.Username = ctx.Query("username")
	userReq.Password = ctx.Query("password")

	userServiceClient := ctx.Keys["user_service"].(service.UserServiceClient)
	userResp, err := userServiceClient.UserLogin(context.Background(), &userReq)
	if err != nil {
		PanicIfUserError(err)
	}
	token, _ := auth.GenerateToken(userResp.UserId)

	r := res.UserResponse{
		StatusCode: userResp.StatusCode,
		StatusMsg:  userResp.StatusMsg,
		UserId:     userResp.UserId,
		Token:      token,
	}

	ctx.JSON(http.StatusOK, r)
}

// UserInfo 用户信息列表
func UserInfo(ctx *gin.Context) {
	var userIds []int64

	// jwt中间件会解析token，然后把user_id放入context中，所以用两种方式都可以获取到user_id
	userIdStr := ctx.Query("user_id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	userIds = append(userIds, userId)

	r := res.UserInfoResponse{
		StatusCode: exception.SUCCESS,
		StatusMsg:  exception.GetMsg(exception.SUCCESS),
		User:       GetUserInfo(userIds, ctx)[0],
	}

	ctx.JSON(http.StatusOK, r)
}

// GetUserInfo 根据用户id，去调取三个服务，拼接出所有的用户信息
func GetUserInfo(userIds []int64, ctx *gin.Context) (userInfos []res.User) {
	var err error
	// 构建三个服务的请求
	var userInfoReq service.UserInfoRequest
	var countInfoReq service.CountRequest
	var followInfoReq service.FollowInfoRequest

	userInfoReq.UserIds = userIds
	countInfoReq.UserIds = userIds
	followInfoReq.ToUserId = userIds

	// 创建接收三个响应
	var userResp *service.UserInfoResponse
	var countInfoResp *service.CountResponse
	var followInfoResp *service.FollowInfoResponse

	// 分别去调用三个服务
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		userServiceClient := ctx.Keys["user_service"].(service.UserServiceClient)
		userResp, err = userServiceClient.UserInfo(context.Background(), &userInfoReq)
		if err != nil {
			PanicIfUserError(err)
		}
	}()

	go func() {
		defer wg.Done()
		videoServiceClient := ctx.Keys["video_service"].(service.VideoServiceClient)
		countInfoResp, err = videoServiceClient.CountInfo(context.Background(), &countInfoReq)
		if err != nil {
			PanicIfVideoError(err)
		}
	}()

	go func() {
		defer wg.Done()

		// 拿到当前用户的id，在is_follow找对应关系
		userIdStr, _ := ctx.Get("user_id")
		userId, _ := userIdStr.(int64)
		followInfoReq.UserId = userId
		socialServiceClient := ctx.Keys["social_service"].(service.SocialServiceClient)
		followInfoResp, err = socialServiceClient.GetFollowInfo(context.Background(), &followInfoReq)
		if err != nil {
			PanicIfFollowError(err)
		}
	}()
	wg.Wait()

	// 构建信息userResp.Users[0], countInfoResp.Counts[0])
	for index, _ := range userIds {
		userInfos = append(userInfos, BuildUser(userResp.Users[index], countInfoResp.Counts[index], followInfoResp.FollowInfo[index]))
	}

	return userInfos
}

// BuildUser 构建用户信息
func BuildUser(user *service.User, count *service.Count, follow *service.FollowInfo) res.User {
	return res.User{
		Id:   user.Id,
		Name: user.Name,

		FollowCount:   follow.FollowCount,
		FollowerCount: follow.FollowerCount,
		IsFollow:      follow.IsFollow,

		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,

		TotalFavorited: strconv.FormatInt(count.TotalFavorited, 10), // 将int64转换成string
		WorkCount:      count.WorkCount,
		FavoriteCount:  count.FavoriteCount,
	}
}
