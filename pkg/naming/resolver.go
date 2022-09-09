package naming

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"github.com/stellarisJAY/goim/pkg/log"
	"google.golang.org/grpc/resolver"
	"net/url"
)

const consulScheme = "consul"

type consulResolver struct {
	cc          resolver.ClientConn // client conn
	serviceName string              // 服务名称
	lastIndex   uint64              // consul 订阅的最后一个index
}

type consulBuilder struct {
}

var (
	WrongSchemeError = errors.New("wrong scheme for resolver")
)

func init() {
	// 注册 resolver 到grpc
	resolver.Register(&consulBuilder{})
}

func (c *consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	_, _, name, err := parseTarget(target.URL)
	if err != nil {
		return nil, err
	}
	cr := &consulResolver{
		cc:          cc,
		lastIndex:   0,
		serviceName: name,
	}
	// 订阅之前先获取当前存在的服务列表
	services, meta, err := client.Health().Service(name, "", true, nil)
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
func (c *consulResolver) watch() {
	for {
		services, meta, err := client.Health().Service(c.serviceName, c.serviceName, true, &api.QueryOptions{WaitIndex: c.lastIndex})
		if err != nil {
			log.Warn("watch healthy services error %v", err)
			break
		}
		addresses := make([]resolver.Address, 0, len(services))
		for _, service := range services {
			log.Info("new address of service %s online : %s", c.serviceName, service.Service.Address)
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

func (c *consulResolver) ResolveNow(options resolver.ResolveNowOptions) {

}

func (c *consulResolver) Close() {

}

func parseTarget(url url.URL) (host, port, serviceName string, err error) {
	if url.Scheme != consulScheme {
		err = WrongSchemeError
		return
	}
	host = url.Host
	port = url.Port()
	serviceName = url.Path
	return
}
