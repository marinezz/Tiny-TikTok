// Package model 数据库自动迁移
package model

func migration() {
	// 自动迁移
	DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})
	// Todo 判断error 写入日志
}
