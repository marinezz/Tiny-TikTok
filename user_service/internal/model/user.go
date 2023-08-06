package model

import "user/internal/service"

type User struct {
	Id              int64  `gorm:"primary_key"`
	UserName        string `gorm:"unique"`
	PassWord        string `gorm:"notnull"`
	Avatar          string // 用户头像
	BackgroundImage string // 用户首页顶部图
	Signature       string `gorm:"default:该用户还没有简介"` // 个人简介
}

// Create 创建用户
func (*User) Create(request *service.UserRegisterRequest) error {
	user := User{
		UserName: request.Username,
		PassWord: request.Password,
	}
	DB.Create(user)
	return nil
}
