package exception

//* 错误码：
//* 四位组成
//* 1. 1开头代表用户端错误
//* 2. 2开头代表当前系统异常
//* 3. 3开头代表第三方服务异常
//* 4. 4开头若无法确定具体错误，选择宏观错误
//* 5. 大的错误类间的步长间距预留100

const (
	SUCCESS = 0  // 常规的返回成功
	ERROR   = -1 // 常规返回失败

	RequestERROR = 1000 // token相关
	UnAuth       = 1001
	TokenTimeOut = 1002

	ErrOperate = 1200 // 异常操作

	DataErr  = 2000 // 数据创建
	CacheErr = 2001 // 缓存异常

	UserExist     = 2100 // 用户相关
	UserUnExist   = 2101
	PasswordError = 2102

	VideoUnExist     = 2200 // 视频相关
	VideoUploadErr   = 2201
	VideoFavoriteErr = 2203
	UserNoVideo      = 2204

	FavoriteErr       = 2300 // 点赞相关
	CancelFavoriteErr = 2301
	VideoNoFavorite   = 2302
	UserNoFavorite    = 2303

	CommentErr       = 2400 // 评论失败
	CommentUnExist   = 2401 // 评论相关
	CommentDeleteErr = 2402

	FollowSelfErr = 2500 // 关注相关
)
