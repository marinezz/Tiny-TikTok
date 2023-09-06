package exception

var Msg = map[int]string{

	SUCCESS: "OK",
	ERROR:   "fail",

	RequestERROR: "请求有误，用户未登录",
	UnAuth:       "Token未经授权",
	TokenTimeOut: "Token授权已过期",

	ErrOperate: "错误操作",

	DataErr:  "数据创建错误",
	CacheErr: "Redis异常",

	UserExist:     "用户已经存在",
	UserUnExist:   "用户不存在，请注册",
	PasswordError: "密码错误",

	VideoUnExist:     "视频信息不存在",
	VideoUploadErr:   "视频上传失败",
	VideoFavoriteErr: "视频点赞失败",
	UserNoVideo:      "用户没有发布作品",

	FavoriteErr:       "点赞失败",
	CancelFavoriteErr: "取消点赞失败",
	VideoNoFavorite:   "视频没有点赞",
	UserNoFavorite:    "用户没有点赞",

	CommentErr:       "评论失败",
	CommentUnExist:   "评论不存在",
	CommentDeleteErr: "评论删除失败",

	FollowSelfErr: "关注自己",
}

// GetMsg 根据状态码获取对应信息
func GetMsg(code int) string {
	msg, ok := Msg[code]
	if ok {
		return msg
	}

	return Msg[ERROR]
}
