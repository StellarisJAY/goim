package naming

import (
	"context"
	"github.com/hashicorp/consul/api"
	"github.com/stellarisJAY/goim/pkg/config"
	"google.golang.org/grpc"
	"log"
)

var consulAddress string
var client *api.Client

type ServiceRegistration struct {
	ID          string
	ServiceName string
	Address     string
}

func init() {
	consulAddress = config.Config.Consul.Address
	// 初始化 consul 客户端
	conf := api.DefaultConfig()
	conf.Address = consulAddress
	c, err := api.NewClient(conf)
	if err != nil {
		panic(err)
	}
	client = c
	log.Println("consul service discovery initialized")
}

// GetClientConn 获取一个指定服务的客户端连接
func GetClientConn(serviceName string) (*grpc.ClientConn, error) {
	return grpc.DialContext(context.Background(), consulScheme+"://"+consulAddress+"/"+serviceName, grpc.WithInsecure())
}

// DialConnection 获取指定地址的客户端连接
func DialConnection(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	return conn, err
}

// RegisterService 注册服务
func RegisterService(registration ServiceRegistration) error {
	err := client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name:    registration.ServiceName,
		Address: registration.Address,
		Kind:    api.ServiceKindTypical,
	})
	if err != nil {
		return err
	}
	log.Printf("service %s registered to consul", registration.ServiceName)
	return nil
}
