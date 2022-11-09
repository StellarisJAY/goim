package naming

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/stellarisJAY/goim/pkg/config"
	"google.golang.org/grpc"
	"strings"
)

type ServiceRegistration struct {
	ID          string
	ServiceName string
	Address     string
}

type Naming interface {
	Init()
	// GetClientConn 获取一个指定服务的客户端连接
	GetClientConn(serviceName string, tracer opentracing.Tracer) (*grpc.ClientConn, error)
	// DialConnection 获取指定地址的客户端连接
	DialConnection(address string) (*grpc.ClientConn, error)
	// RegisterService 注册服务
	RegisterService(registration ServiceRegistration) error
	// CurrentServiceName 当前服务器的service名称
	CurrentServiceName() string
}

var ns Naming

func init() {
	n := strings.ToLower(config.Config.Naming)
	switch n {
	case "consul":
		ns = &ConsulNaming{}
	case "etcd":
		ns = &EtcdNaming{}
	default:
		panic(fmt.Errorf("unknown or unsupported naming system: %s", n))
	}
	ns.Init()
}

// GetClientConn 获取一个指定服务的客户端连接
func GetClientConn(serviceName string, tracer opentracing.Tracer) (*grpc.ClientConn, error) {
	return ns.GetClientConn(serviceName, tracer)
}

// DialConnection 获取指定地址的客户端连接
func DialConnection(address string) (*grpc.ClientConn, error) {
	return ns.DialConnection(address)
}

// RegisterService 注册服务
func RegisterService(registration ServiceRegistration) error {
	return ns.RegisterService(registration)
}

func CurrentServiceName() string {
	return ns.CurrentServiceName()
}
