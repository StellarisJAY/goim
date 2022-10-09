package naming

import (
	"context"
	"github.com/hashicorp/consul/api"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

type ConsulNaming struct {
	consulAddress string
	client        *api.Client
}

func (ns *ConsulNaming) Init() {
	ns.consulAddress = config.Config.Consul.Address
	// 初始化 consul 客户端
	conf := api.DefaultConfig()
	conf.Address = ns.consulAddress
	c, err := api.NewClient(conf)
	if err != nil {
		panic(err)
	}
	ns.client = c
	// 注册 resolver 到grpc
	resolver.Register(&consulBuilder{ns: ns})
	log.Info("consul service discovery initialized")
}

// GetClientConn 获取一个指定服务的客户端连接
func (ns *ConsulNaming) GetClientConn(serviceName string) (*grpc.ClientConn, error) {
	return grpc.DialContext(context.Background(), consulScheme+"://"+ns.consulAddress+"/"+serviceName, grpc.WithInsecure())
}

// DialConnection 获取指定地址的客户端连接
func (ns *ConsulNaming) DialConnection(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	return conn, err
}

// RegisterService 注册服务
func (ns *ConsulNaming) RegisterService(registration ServiceRegistration) error {
	err := ns.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name:    registration.ServiceName,
		Address: registration.Address,
		Kind:    api.ServiceKindTypical,
	})
	if err != nil {
		return err
	}
	return nil
}
