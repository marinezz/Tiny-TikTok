package exception

var Msg = map[int]string{
	SUCCESS: "OK",
	ERROR:   "fail",

	RequestERROR: "请求有误",
	UnAuth:       "Token未经授权",
	TokenTimeOut: "Token授权已过期",

	UserExist: "用户已经存在",
}

// GetMsg 根据状态码获取对应信息
func GetMsg(code int) string {
	msg, ok := Msg[code]
	if ok {
		return msg
	}

	return Msg[ERROR]
}
