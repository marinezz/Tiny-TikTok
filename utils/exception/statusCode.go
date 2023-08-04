package exception

//* 错误码：
//* 四位组成
//* 1. 1开头代表用户端错误
//* 2. 2开头代表当前系统异常
//* 3. 3开头代表第三方服务异常
//* 4. 4开头若无法确定具体错误，选择宏观错误
//* 5. 大的错误类间的步长间距预留100

const (
	SUCCESS = 200
	ERROR   = 500

	RequestERROR = 1000
	UnAuth       = 1001
	TokenTimeOut = 1002
)
