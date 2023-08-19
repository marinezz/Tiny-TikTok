package handler

import (
	"github.com/pkg/errors"
)

// PanicIfVideoError 视频错误处理
func PanicIfVideoError(err error) {
	if err != nil {
		err = errors.New("videoService--error--" + err.Error())
		// Todo 统一的日志处理
		panic(err)
	}
}

// PanicIfUserError 用户错误处理
func PanicIfUserError(err error) {
	if err != nil {
		err = errors.New("UserService--error" + err.Error())
		panic(err)
	}
}

// PanicIfMessageError 消息错误处理
func PanicIfMessageError(err error) {
	if err != nil {
		err = errors.New("MessageService--error" + err.Error())
		panic(err)
	}
}

// PanicIfFollowError 关注错误处理
func PanicIfFollowError(err error) {
	if err != nil {
		err = errors.New("FollowService--error" + err.Error())
		panic(err)
	}
}

// PanicIfFavoriteError 喜欢错误处理
func PanicIfFavoriteError(err error) {
	if err != nil {
		err = errors.New("FavoriteService--error" + err.Error())
		panic(err)
	}
}

// PanicIfCommentError 评论错误处理
func PanicIfCommentError(err error) {
	if err != nil {
		err = errors.New("CommentService--error" + err.Error())
		panic(err)
	}
}
