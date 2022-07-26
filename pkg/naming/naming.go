package naming

import (
	"context"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"log"
)

var consulAddress = "127.0.0.1:8500"
var client *api.Client

type ServiceRegistration struct {
	ID          string
	ServiceName string
	Address     string
}

func init() {
	// 初始化 consul 客户端
	config := api.DefaultConfig()
	config.Address = consulAddress
	c, err := api.NewClient(config)
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
