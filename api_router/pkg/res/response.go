package res

type Comment struct {
	Content    string `json:"content"`     // 评论内容
	CreateDate string `json:"create_date"` // 评论发布日期，格式 mm-dd
	ID         int64  `json:"id"`          // 评论id
	// Todo 用户
}

type UserResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
}

type UserInfoResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg,omitempty"`
	User       interface{} `json:"user"`
}

type PublishActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type FavoriteActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type CommentActionResponse struct {
	StatusCode int32    `json:"status_code"`
	StatusMsg  string   `json:"status_msg"`
	Comment    *Comment `json:"comment"` // 评论成功返回评论内容，不需要重新拉取整个列表
}

type FollowActionResponse struct {
	StatusCode int32  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

type FollowListResponse struct {
	StatusCode int32   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string  `json:"status_msg"`  // 返回状态描述
	UserList   []int64 `json:"user_list"`
}

type PostMessageResponse struct {
	StatusCode int32  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

type Message struct {
	Id         int64  `json:"id"`
	ToUserId   int64  `json:"to_user_id"`
	FromUserID int64  `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime string `json:"create_time"`
}

type GetMessageResponse struct {
	StatusCode  int32     `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg   string    `json:"status_msg"`  // 返回状态描述
	MessageList []Message `json:"message_list"`
}
