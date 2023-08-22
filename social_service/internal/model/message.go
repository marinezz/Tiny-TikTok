package model

import (
	"sync"
	"time"
	"utils/snowFlake"
)

type Message struct {
	Id        int64 `gorm:"primary_key"`
	UserId    int64
	ToUserId  int64
	Message   string
	CreatedAt string
}

type MessageModel struct {
}

var messageModel *MessageModel
var messageOnce sync.Once // 单例模式

// GetMessageInstance 获取单例实例
func GetMessageInstance() *MessageModel {
	messageOnce.Do(
		func() {
			messageModel = &MessageModel{}
		},
	)
	return messageModel
}

func getCurrentTime() (createTime string) {
	createTime = time.Now().Format("2006-01-02 15:04:05")
	return
}

func (*MessageModel) PostMessage(message *Message) error {
	message.CreatedAt = getCurrentTime()
	flake, _ := snowFlake.NewSnowFlake(7, 3)
	message.Id = flake.NextId()
	err := DB.Create(&message).Error
	return err
}

func (*MessageModel) GetMessage(UserId int64, ToUserID int64, messages *[]Message) error {
	err := DB.Model(&Message{}).Where(&Message{UserId: UserId, ToUserId: ToUserID}).Order("created_at").Find(messages).Error
	return err
}
