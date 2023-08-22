package wrapper

import (
	"api_router/pkg/logger"
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
)

func NewWrapper(name string) {
	// 设置 Hystrix 配置
	hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
		Timeout:                1000, // 超时时间（毫秒）
		MaxConcurrentRequests:  10,   // 最大并发请求数
		RequestVolumeThreshold: 5,    // 触发熔断的最小请求数，大于这个值才开始做熔断检测
		SleepWindow:            5000, // 熔断后休眠时间（毫秒）
		ErrorPercentThreshold:  20,   // 错误率阈值
	})

	// 使用 Hystrix 执行命令
	err := hystrix.Do(name, func() error {
		return errors.New("服务熔断")
	}, nil)
	if err != nil {
		logger.Log.Info("请求过快---" + name + "服务熔断")
		fmt.Println("err", err)
	}

}
