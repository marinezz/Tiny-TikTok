package etcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

// EtcdRegister 服务注册
type EtcdRegister struct {
	etcdCli *clientv3.Client // etcdClient对象
	leaseId clientv3.LeaseID // 租约id
}

// CreateLease 创建租约。expire表示有效期(s)
func (e *EtcdRegister) CreateLease(expire int64) error {

	lease, err := e.etcdCli.Grant(context.Background(), expire)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	e.leaseId = lease.ID // 记录生成的租约Id
	return nil
}

// BindLease 绑定租约。将租约与对应的key-value绑定
func (e *EtcdRegister) BindLease(key string, value string) error {

	res, err := e.etcdCli.Put(context.Background(), key, value, clientv3.WithLease(e.leaseId))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Printf("bind lease success %v \n", res)
	return nil
}

// KeepAlive 获取续约通道 并 持续续租
func (e *EtcdRegister) KeepAlive() error {
	keepAliveChan, err := e.etcdCli.KeepAlive(context.Background(), e.leaseId)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	go func(keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse) {
		for {
			select {
			case resp := <-keepAliveChan:
				log.Printf("续约成功...leaseID=%d", resp.ID)
			}
		}
	}(keepAliveChan)

	return nil
}

// Close 关闭服务
func (e *EtcdRegister) Close() error {
	log.Printf("close...\n")
	// 撤销租约
	e.etcdCli.Revoke(context.Background(), e.leaseId)
	return e.etcdCli.Close()
}

// NewEtcdRegister 初始化etcd服务注册对象
func NewEtcdRegister(etcdServerAddr string) (*EtcdRegister, error) {

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdServerAddr},
		DialTimeout: 3 * time.Second,
	})

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	e := &EtcdRegister{
		etcdCli: client,
	}
	return e, nil
}

// ServiceRegister 服务注册。expire表示过期时间,serviceName和serviceAddr分别是服务名与服务地址
func (e *EtcdRegister) ServiceRegister(serviceName, serviceAddr string, expire int64) (err error) {

	// 创建租约
	err = e.CreateLease(expire)
	if err != nil {
		return err
	}

	// 将租约与k-v绑定
	err = e.BindLease(serviceName, serviceAddr)
	if err != nil {
		return err
	}

	// 持续续租
	err = e.KeepAlive()
	return err
}
