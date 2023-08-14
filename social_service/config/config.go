package config

import (
	"github.com/spf13/viper"
	"path"
	"runtime"
	"strings"
)

// Config 数据结构
type Config struct {
	Mysql struct {
		Host     string
		Port     string
		Username string
		Password string
		Database string
	}
}

// ConfigData 配置数据变量
var ConfigData *Config

// InitConfig 读取配置文件
func InitConfig() {
	_, filePath, _, _ := runtime.Caller(0)

	currentDir := path.Dir(filePath)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(currentDir)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

}

// DbDnsInit 拼接链接数据库的DNS
func DbDnsInit() string {
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	username := viper.GetString("mysql.username")
	password := viper.GetString("mysql.password")
	database := viper.GetString("mysql.database")

	InitConfig()
	dns := strings.Join([]string{username, ":", password, "@tcp(", host, ":", port, ")/", database, "?charset=utf8&parseTime=True&loc=Local"}, "")

	return dns
}
