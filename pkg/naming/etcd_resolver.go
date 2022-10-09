package naming

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc/resolver"
	"strings"
)

const etcdScheme = "etcd"
const serviceNamePrefix = "/goim/services/"

type EtcdResolver struct {
	ns *EtcdNaming
	cc resolver.ClientConn
}

type EtcdResolverBuilder struct {
	ns *EtcdNaming
}

func (r *EtcdResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	name, err := parseEtcdTarget(target)
	if err != nil {
		return nil, err
	}
	prefix := serviceNamePrefix + name + "/"
	result, err := r.ns.client.Get(context.TODO(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("etcd get kv error: %w", err)
	}

	rs := &EtcdResolver{cc: cc, ns: r.ns}
	addrList := make([]resolver.Address, len(result.Kvs))
	for i, kv := range result.Kvs {
		addr := strings.TrimPrefix(string(kv.Key), prefix)
		addrList[i] = resolver.Address{Addr: addr}
	}
	cc.NewAddress(addrList)
	go rs.watch(prefix, addrList)
	return rs, nil
}

func (r *EtcdResolver) watch(prefix string, addrList []resolver.Address) {
	watchChan := r.ns.client.Watch(context.TODO(), prefix, clientv3.WithPrefix())
	for response := range watchChan {
		events := response.Events
		for _, event := range events {
			addr := strings.TrimPrefix(string(event.Kv.Key), prefix)
			switch event.Type {
			case clientv3.EventTypePut:
				if !exists(addrList, addr) {
					addrList = append(addrList, resolver.Address{Addr: addr})
				}
			case clientv3.EventTypeDelete:
				for i := 0; i < len(addrList); i++ {
					if addrList[i].Addr == addr {
						addrList[i] = addrList[len(addrList)-1]
						addrList = addrList[:len(addrList)-1]
						break
					}
				}
			default:
				continue
			}
		}
		r.cc.UpdateState(resolver.State{Addresses: addrList})
	}
}

func exists(addrList []resolver.Address, addr string) bool {
	for _, address := range addrList {
		if address.Addr == addr {
			return true
		}
	}
	return false
}

func (r *EtcdResolverBuilder) Scheme() string {
	return etcdScheme
}

func (r *EtcdResolver) ResolveNow(options resolver.ResolveNowOptions) {
}

func (r *EtcdResolver) Close() {
}

func parseEtcdTarget(target resolver.Target) (serviceName string, err error) {
	if target.Scheme != consulScheme {
		err = fmt.Errorf("invalid url scheme")
		return
	}
	serviceName = target.Endpoint
	return
}
