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

	UserExist     = 1100 // 数据库相关
	UserUnExist   = 1101
	PasswordError = 1102

	ErrOperate = 1200 // 异常操作
)
