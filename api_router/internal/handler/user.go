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
)

// UserRegister 用户注册
func UserRegister(ctx *gin.Context) {
	var userReq service.UserRequest
	ctx.Bind(&userReq)

	userServiceClient := ctx.Keys["user_service"].(service.UserServiceClient)
	userResp, err := userServiceClient.UserRegister(context.Background(), &userReq)
	if err != nil {
		panic(err)
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

func UserLogin(ctx *gin.Context) {
	var userReq service.UserRequest
	ctx.Bind(&userReq)

	userServiceClient := ctx.Keys["user_service"].(service.UserServiceClient)
	userResp, err := userServiceClient.UserLogin(context.Background(), &userReq)
	if err != nil {
		panic(err)
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

func UserInfo(ctx *gin.Context) {
	var userInfoReq service.UserInfoRequest
	var countInfoReq service.CountRequest

	// jwt中间件会解析token，然后把user_id放入context中，所以用两种方式都可以获取到user_id
	userIdStr := ctx.Query("user_id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	userInfoReq.UserIds = append(userInfoReq.UserIds, userId)
	countInfoReq.UserIds = append(countInfoReq.UserIds, userId)
	//userId, _ := ctx.Get("user_id")
	//userInfoReq.UserId, _ = userId.(int64)

	userServiceClient := ctx.Keys["user_service"].(service.UserServiceClient)
	userResp, _ := userServiceClient.UserInfo(context.Background(), &userInfoReq)

	videoServiceClient := ctx.Keys["video_service"].(service.VideoServiceClient)
	countInfoResp, _ := videoServiceClient.CountInfo(context.Background(), &countInfoReq)

	r := res.UserInfoResponse{
		StatusCode: userResp.StatusCode,
		StatusMsg:  userResp.StatusMsg,
		User:       BuildUser(userResp.Users[0], countInfoResp.Counts[0]),
	}
	ctx.JSON(http.StatusOK, r)

}

// GetUserInfo 根据用户id，去调取三个服务，拼接出所有的用户信息
func GetUserInfo(userIds []int64, ctx *gin.Context) (userInfos []res.User) {
	// 构建三个服务的请求
	// todo 社交服务
	var userInfoReq service.UserInfoRequest
	var countInfoReq service.CountRequest

	userInfoReq.UserIds = userIds
	countInfoReq.UserIds = userIds

	// 创建接收三个响应
	var userResp *service.UserInfoResponse
	var countInfoResp *service.CountResponse

	// 分别去调用三个服务
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		userServiceClient := ctx.Keys["user_service"].(service.UserServiceClient)
		userResp, _ = userServiceClient.UserInfo(context.Background(), &userInfoReq)
	}()

	go func() {
		defer wg.Done()
		videoServiceClient := ctx.Keys["video_service"].(service.VideoServiceClient)
		countInfoResp, _ = videoServiceClient.CountInfo(context.Background(), &countInfoReq)
	}()

	go func() {
		defer wg.Done()
		// todo 社交服务
	}()
	wg.Wait()

	// 构建信息userResp.Users[0], countInfoResp.Counts[0])
	for index, _ := range userIds {
		userInfos = append(userInfos, BuildUser(userResp.Users[index], countInfoResp.Counts[index]))
	}

	return userInfos
}

// BuildUser 构建用户信息 Todo 还有其余信息的构建
func BuildUser(user *service.User, count *service.Count) res.User {
	return res.User{
		Id:   user.Id,
		Name: user.Name,

		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,

		TotalFavorited: strconv.FormatInt(count.TotalFavorited, 10), // 将int64转换成string
		WorkCount:      count.WorkCount,
		FavoriteCount:  count.FavoriteCount,
	}
}
