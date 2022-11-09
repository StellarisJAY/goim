package naming

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/opentracing/opentracing-go"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"time"
)

type EtcdNaming struct {
	client       *clientv3.Client
	registration ServiceRegistration
}

func (ns *EtcdNaming) Init() {
	endpoints := config.Config.Etcd.Endpoints
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(fmt.Errorf("create etcd client error %w", err))
	}
	ns.client = client
	resolver.Register(&EtcdResolverBuilder{ns: ns})
}

func (ns *EtcdNaming) GetClientConn(serviceName string, tracer opentracing.Tracer) (*grpc.ClientConn, error) {
	options := buildDialOptions(tracer)
	return grpc.DialContext(context.TODO(), etcdScheme+"://host:8888/"+serviceName, options...)
}

func (ns *EtcdNaming) DialConnection(address string) (*grpc.ClientConn, error) {
	options := buildDialOptions(opentracing.GlobalTracer())
	return grpc.Dial(address, options...)
}

func (ns *EtcdNaming) CurrentServiceName() string {
	return ns.registration.ServiceName
}

func (ns *EtcdNaming) RegisterService(registration ServiceRegistration) error {
	ns.registration = registration
	lease, err := ns.client.Grant(context.TODO(), 5000)
	if err != nil {
		return fmt.Errorf("etcd grant lease error: %w", err)
	}
	_, err = ns.client.Put(context.TODO(), serviceNamePrefix+registration.ServiceName+"/"+registration.Address, "", clientv3.WithLease(lease.ID))
	if err != nil {
		return fmt.Errorf("register error, etcd put key value failed: %w", err)
	}
	if err = ns.keepAlive(lease.ID, registration.ServiceName); err != nil {
		return fmt.Errorf("register error, etcd keepalive failed: %w", err)
	}
	return nil
}

func (ns *EtcdNaming) keepAlive(leaseID clientv3.LeaseID, serviceName string) error {
	aliveChan, err := ns.client.KeepAlive(context.TODO(), leaseID)
	if err != nil {
		return err
	}
	go func(serviceName string) {
		for {
			select {
			case _, ok := <-aliveChan:
				if !ok {
					log.Warn("service register keepalive interrupted", zap.String("serviceName", serviceName))
					break
				}
			}
		}
	}(serviceName)
	return nil
}
