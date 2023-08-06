package model

import (
	"testing"
	"user/internal/service"
)

func TestUser_Create(t *testing.T) {
	InitDb()
	u := new(User)
	req := new(service.UserRegisterRequest)
	req.Username = "marine"
	req.Password = "123456"

	u.Create(req)

}
