package naming

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/stellarisJAY/goim/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
)

const consulScheme = "consul"

type ConsulResolver struct {
	cc          resolver.ClientConn // client conn
	serviceName string              // 服务名称
	lastIndex   uint64              // consul 订阅的最后一个index
	ns          *ConsulNaming
}

type consulBuilder struct {
	ns *ConsulNaming
}

func init() {

}

func (c *consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	name, err := parseConsulTarget(target)
	if err != nil {
		return nil, err
	}
	cr := &ConsulResolver{
		cc:          cc,
		lastIndex:   0,
		serviceName: name,
		ns:          c.ns,
	}
	// 订阅之前先获取当前存在的服务列表
	services, meta, err := c.ns.client.Health().Service(name, "", true, nil)
	if err != nil {
		return nil, err
	}
	addresses := make([]resolver.Address, 0, len(services))
	for _, service := range services {
		addresses = append(addresses, resolver.Address{Addr: service.Service.Address})
	}
	cc.NewAddress(addresses)
	cc.NewServiceConfig(name)
	cr.lastIndex = meta.LastIndex
	// 开启watch订阅
	go cr.watch()
	return cr, nil
}

// watch 订阅 consul 服务列表的变化
func (c *ConsulResolver) watch() {
	for {
		services, meta, err := c.ns.client.Health().Service(c.serviceName, c.serviceName, true, &api.QueryOptions{WaitIndex: c.lastIndex})
		if err != nil {
			log.Warn("retrieve healthy services failed",
				zap.String("serviceName", c.serviceName),
				zap.Error(err))
			break
		}
		addresses := make([]resolver.Address, 0, len(services))
		for _, service := range services {
			log.Info("new service instance online",
				zap.String("serviceName", c.serviceName),
				zap.String("address", service.Service.Address))
			addresses = append(addresses, resolver.Address{Addr: service.Service.Address})
		}
		c.lastIndex = meta.LastIndex
		c.cc.NewAddress(addresses)
		c.cc.NewServiceConfig(c.serviceName)
	}
}

func (c *consulBuilder) Scheme() string {
	return consulScheme
}

func (c *ConsulResolver) ResolveNow(options resolver.ResolveNowOptions) {

}

func (c *ConsulResolver) Close() {

}

func parseConsulTarget(target resolver.Target) (serviceName string, err error) {
	if target.Scheme != consulScheme {
		err = fmt.Errorf("invalid url scheme")
		return
	}
	serviceName = target.Endpoint
	return
}
