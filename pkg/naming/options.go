package naming

import (
	"github.com/opentracing/opentracing-go"
	"github.com/stellarisJAY/goim/pkg/trace"
	"google.golang.org/grpc"
)

func buildDialOptions(tracer opentracing.Tracer) []grpc.DialOption {
	options := []grpc.DialOption{grpc.WithInsecure()}
	if tracer != nil {
		clientInterceptor := trace.ClientInterceptor(tracer)
		options = append(options, grpc.WithUnaryInterceptor(clientInterceptor))
	}
	return options
}
