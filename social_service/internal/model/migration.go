// Package model 数据库自动迁移
package model

import "log"

func migration() {
	// 自动迁移
	err := DB.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(&Follow{})
	// Todo 判断error 写入日志
	if err != nil {
		log.Print("err")
	}

	err = DB.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(&Message{})
	// Todo 判断error 写入日志
	if err != nil {
		log.Print("err")
	}
}
