package wrapper

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"time"
)

// NewServiceWrapper todo 服务熔断改一改
func NewServiceWrapper(name string) {
	c := &Config{
		Namespace:              name,
		Timeout:                1 * time.Second, // TODO 建议加在配置文件里面
		MaxConcurrentRequests:  3,
		RequestVolumeThreshold: 5,
		SleepWindow:            3 * time.Second,
		ErrorPercentThreshold:  50,
	}

	g := NewGroup(c)

	if err := g.Do(name, func() error {
		log.Print("请求太快，请稍后再试...")
		return errors.New("服务熔断")
	}); err != nil {
		log.Print("请求太快，请稍后再试...")
		fmt.Println("err", err)
	}
}
