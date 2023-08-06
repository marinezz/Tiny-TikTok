package config

import "github.com/spf13/viper"

// InitConfig 读取配置文件
func InitConfig() {

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("conf")

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

}
