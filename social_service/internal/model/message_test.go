package model

import (
	"fmt"
	"testing"
)

func TestMessageModel_GetMessage(t *testing.T) {
	InitDb()
	var messages []Message
	err := GetMessageInstance().GetMessage(1, 2, &messages)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", messages)
}

func TestMessageModel_PostMessage(t *testing.T) {
	InitDb()
	err := GetMessageInstance().PostMessage(&Message{
		UserId:   1,
		ToUserId: 2,
		Message:  "你好",
	})
	if err != nil {
		panic(err)
	}
}
