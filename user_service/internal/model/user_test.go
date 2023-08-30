package model

import (
	"fmt"
	"testing"
)

func TestUser_Create(t *testing.T) {
	InitDb()
	user := &User{
		UserName: "张三",
		PassWord: "123456",
	}
	GetInstance().Create(user)
}

func TestUserModel_FindUserByName(t *testing.T) {
	InitDb()
	user, _ := GetInstance().FindUserByName("ben")
	fmt.Print(user)
}

func TestUserModel_CheckUserExist(t *testing.T) {
	InitDb()
	exist := GetInstance().CheckUserExist("lisi")
	fmt.Println(exist)
}

func TestUserModel_FindUserById(t *testing.T) {
	InitDb()
	user, _ := GetInstance().FindUserById(4391185549234176)
	fmt.Println(user)
}
