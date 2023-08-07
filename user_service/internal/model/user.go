package model

import (
	"sync"
	"user/pkg/encryption"
	"utils/snowFlake"
)

type User struct {
	Id              int64  `gorm:"primary_key"`
	UserName        string `gorm:"unique"`
	PassWord        string `gorm:"notnull"`
	Avatar          string // 用户头像
	BackgroundImage string // 用户首页顶部图
	Signature       string `gorm:"default:该用户还没有简介"` // 个人简介
}

type UserModel struct {
}

var userModel *UserModel
var userOnce sync.Once // 单例模式

// GetInstance 获取单例实例
func GetInstance() *UserModel {
	userOnce.Do(
		func() {
			userModel = &UserModel{}
		},
	)
	return userModel
}

// Create 创建用户
func (*UserModel) Create(user *User) error {
	flake, _ := snowFlake.NewSnowFlake(7, 1)
	user.Id = flake.NextId()
	user.PassWord = encryption.HashPassword(user.PassWord)
	DB.Create(&user)
	return nil
}

// FindUserByName 根据用户名称查找用户,并返回对象
func (*UserModel) FindUserByName(username string) (*User, error) {
	user := User{}
	res := DB.Where("user_name=?", username).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}
