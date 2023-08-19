// 服务发现,发现所有的服务，返回一个map

package discovery

import (
	"api_router/internal/service"
	"api_router/pkg/wrapper"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"utils/etcd"
)

func Resolver() map[string]interface{} {
	serveInstance := make(map[string]interface{})

	etcdAddress := viper.GetString("etcd.address")
	serviceDiscovery, err := etcd.NewServiceDiscovery([]string{etcdAddress})
	if err != nil {
		log.Fatal(err)
	}
	defer serviceDiscovery.Close()

	// 获取用户服务实例
	err = serviceDiscovery.ServiceDiscovery("user_service")
	if err != nil {
		log.Fatal(err)
	}
	userServiceAddr, _ := serviceDiscovery.GetService("user_service")
	userConn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	userClient := service.NewUserServiceClient(userConn)
	serveInstance["user_service"] = userClient

	// 获取视频服务实例
	err = serviceDiscovery.ServiceDiscovery("video_service")
	if err != nil {
		log.Fatal(err)
	}
	videoServiceAddr, _ := serviceDiscovery.GetService("video_service")
	videoConn, err := grpc.Dial(videoServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	videoClient := service.NewVideoServiceClient(videoConn)
	serveInstance["video_service"] = videoClient

	// 获取社交服务实例
	err = serviceDiscovery.ServiceDiscovery("social_service")
	if err != nil {
		log.Fatal(err)
	}
	socialServiceAddr, _ := serviceDiscovery.GetService("social_service")
	socialConn, err := grpc.Dial(socialServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	socialClient := service.NewSocialServiceClient(socialConn)
	serveInstance["social_service"] = socialClient

	wrapper.NewServiceWrapper("user_service")
	wrapper.NewServiceWrapper("video_service")
	wrapper.NewServiceWrapper("social_service")

	return serveInstance
}
