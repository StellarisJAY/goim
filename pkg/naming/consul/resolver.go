package consul

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"github.com/stellarisJAY/goim/pkg/log"
	"google.golang.org/grpc/resolver"
	"net/url"
)

const consulScheme = "consul"

type Resolver struct {
	cc          resolver.ClientConn // client conn
	serviceName string              // 服务名称
	lastIndex   uint64              // consul 订阅的最后一个index
	ns          *Naming
}

type consulBuilder struct {
	ns *Naming
}

var (
	WrongSchemeError = errors.New("wrong scheme for resolver")
)

func init() {

}

func (c *consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	_, _, name, err := parseTarget(target.URL)
	if err != nil {
		return nil, err
	}
	cr := &Resolver{
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
func (c *Resolver) watch() {
	for {
		services, meta, err := c.ns.client.Health().Service(c.serviceName, c.serviceName, true, &api.QueryOptions{WaitIndex: c.lastIndex})
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

func (c *Resolver) ResolveNow(options resolver.ResolveNowOptions) {

}

func (c *Resolver) Close() {

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
