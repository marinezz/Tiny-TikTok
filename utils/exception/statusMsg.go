package exception

var Msg = map[int]string{

	SUCCESS: "OK",
	ERROR:   "fail",

	RequestERROR: "请求有误，用户未登录",
	UnAuth:       "Token未经授权",
	TokenTimeOut: "Token授权已过期",

	UserExist:     "用户已经存在",
	UserUnExist:   "用户不存在，请注册",
	PasswordError: "密码错误",

	ErrOperate: "错误操作",
}

// GetMsg 根据状态码获取对应信息
func GetMsg(code int) string {
	msg, ok := Msg[code]
	if ok {
		return msg
	}

	return Msg[ERROR]
}
