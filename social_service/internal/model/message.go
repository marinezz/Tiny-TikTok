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
	CreatedAt int64
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

func getCurrentTime() (createTime int64) {
	createTime = time.Now().Unix()
	return
}

func (*MessageModel) PostMessage(message *Message) error {
	message.CreatedAt = getCurrentTime()
	flake, _ := snowFlake.NewSnowFlake(7, 3)
	message.Id = flake.NextId()
	err := DB.Create(&message).Error
	return err
}

func (*MessageModel) GetMessage(UserId int64, ToUserID int64, PreMsgTime int64, messages *[]Message) error {
	err := DB.Model(&Message{}).Where(DB.Model(&Message{}).Where(&Message{UserId: UserId, ToUserId: ToUserID}).
		Or(&Message{UserId: ToUserID, ToUserId: UserId})).Where("created_at > ?", PreMsgTime).
		Order("created_at").Find(messages).Error
	return err
}
