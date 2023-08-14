package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func FollowAction(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
