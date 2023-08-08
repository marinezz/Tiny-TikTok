// 服务发现,发现所有的服务，并写入context中

package discovery

import (
	"api_router/internal/service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"utils/etcd"
)

func resolver() {
	etcdAddress := viper.GetString("etcd.address")
	serviceDiscovery, err := etcd.NewServiceDiscovery([]string{etcdAddress})
	if err != nil {
		log.Fatal(err)
	}
	defer serviceDiscovery.Close()

	err = serviceDiscovery.ServiceDiscovery("user_service")
	if err != nil {
		log.Fatal(err)
	}

	serviceAddr, _ := serviceDiscovery.GetService("user_service")
	conn, err := grpc.Dial(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := service.NewUserServiceClient(conn)

	// Todo 将服务实例放入context中
}
