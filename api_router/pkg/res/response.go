package res

type User struct {
	// 用户服务
	Id   int64  `json:"id"`
	Name string `json:"name"`
	// 社交服务
	FollowCount   int64 `json:"follow_count"`
	FollowerCount int64 `json:"follower_count"`
	IsFollow      bool  `json:"is_follow"`
	// 用户服务
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	// 视频服务
	TotalFavorited string `json:"total_favorited"` // 获赞数量
	WorkCount      int64  `json:"work_count"`      // 作品数量
	FavoriteCount  int64  `json:"favorite_count"`  // 喜欢数量
}

type Comment struct {
	Id         int64  `json:"id"` // 评论id
	User       User   `json:"user"`
	Content    string `json:"content"`     // 评论内容
	CreateDate string `json:"create_date"` // 评论发布日期，格式 mm-dd
}

type Video struct {
	Id            int64  `json:"id"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}

type FeedResponse struct {
	StatusCode int64   `json:"status_code"`
	NextTime   int64   `json:"next_time"`
	StatusMsg  string  `json:"status_msg"`
	VideoList  []Video `json:"video_list"`
}

type UserResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
}

type UserInfoResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
	User       User   `json:"user"`
}

type VideoListResponse struct {
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string  `json:"status_msg"`  // 返回状态描述
	VideoList  []Video `json:"video_list"`  // 用户点赞视频列表
}

type PublishActionResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type FavoriteActionResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type CommentActionResponse struct {
	StatusCode int64   `json:"status_code"`
	StatusMsg  string  `json:"status_msg"`
	Comment    Comment `json:"comment"` // 评论成功返回评论内容，不需要重新拉取整个列表
}

type CommentDeleteResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type CommentListResponse struct {
	StatusCode int64     `json:"status_code"`
	StatusMsg  string    `json:"status_msg"`
	Comments   []Comment `json:"comment_list"`
}

type FollowActionResponse struct {
	StatusCode int32  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

type FollowListResponse struct {
	StatusCode int32  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	UserList   []User `json:"user_list"`
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
	CreateTime int64  `json:"create_time"`
}

type GetMessageResponse struct {
	StatusCode  int32     `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg   string    `json:"status_msg"`  // 返回状态描述
	MessageList []Message `json:"message_list"`
}
