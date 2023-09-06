package discovery

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"utils/etcd"
	"video/internal/handler"
	"video/internal/service"
)

// AutoRegister etcd自动注册
func AutoRegister() {
	etcdAddress := viper.GetString("etcd.address")
	etcdRegister, err := etcd.NewEtcdRegister(etcdAddress)

	if err != nil {
		log.Fatal(err)
	}

	serviceName := viper.GetString("server.name")
	serviceAddress := viper.GetString("server.address")
	err = etcdRegister.ServiceRegister(serviceName, serviceAddress, 30)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		log.Fatal(err)
	}

	// 传输限制，默认4MB
	var options = []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024 * 128),
	}
	server := grpc.NewServer(options...)
	service.RegisterVideoServiceServer(server, handler.NewVideoService())

	err = server.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
