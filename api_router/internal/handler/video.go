package handler

import (
	"api_router/internal/service"
	"github.com/gin-gonic/gin"
)

func PublishAction(ctx *gin.Context) {
	var publishActionReq service.PublishActionRequest

	userId, _ := ctx.Get("user_id")
	publishActionReq.UserId = userId.(int64)

}
