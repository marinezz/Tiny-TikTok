package model

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
	"user/configs/config"
)

// DB 创建数据库单例
var DB *gorm.DB

// InitDb 初始化数据库
func InitDb() {
	config.InitConfig()
	dns := config.DbDnsInit()
	Database(dns)
}

// Database 初始花数据库链接
func Database(connString string) {
	var ormLogger logger.Interface

	// gin中有debug 、 release 、 test 三种模式
	// 不指定默认以debug形式启动
	// 开发用debug模式、 上线用release模式
	// 指定方式 : gin.SetMode(gin.ReleaseMode)

	// 设置gorm的日志模式，可以打印原生SQL语句
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       connString, // DSN data source name
		DefaultStringSize:         256,        // string 类型字段的默认长度
		DisableDatetimePrecision:  true,       // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,       // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,       // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,      // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 默认不加负数
		},
	})

	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()

	// 设置连接池
	sqlDB.SetMaxIdleConns(20)                  // 空闲时候的最大连接数
	sqlDB.SetMaxOpenConns(100)                 // 打开时候的最大连接数
	sqlDB.SetConnMaxLifetime(time.Second * 20) // 超时时间
	DB = db

	migration()
}
