package trace

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	_log "github.com/stellarisJAY/goim/pkg/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
)

// ClientInterceptor 客户端的trace拦截器
func ClientInterceptor(tracer opentracing.Tracer) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		span, _ := opentracing.StartSpanFromContext(ctx, "gRPC Call", nil)
		defer span.Finish()
		// 从context获得grpc的metadata
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			// 拷贝metadata避免直接在原来的md上面操作
			md = md.Copy()
		}
		// 将tracer信息注入metadata
		err := tracer.Inject(span.Context(), opentracing.TextMap, &MetaDataReadWriter{md})
		if err != nil {
			span.LogFields(log.String("inject-error", err.Error()))
		}
		// 使用invoker发送rpc请求
		outgoingCtx := metadata.NewOutgoingContext(ctx, md)
		err = invoker(outgoingCtx, method, req, reply, cc, opts...)
		if err != nil {
			span.LogFields(log.String("rpc-call-error", err.Error()))
		}
		return err
	}
}

// ServerInterceptor 服务端trace拦截器
func ServerInterceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		md := metadata.New(nil)
		// 从context提取metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}
		// 从metadata提取spanContext
		spanContext, err := tracer.Extract(opentracing.TextMap, &MetaDataReadWriter{md})
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			_log.Warn("extract span from metadata error", zap.Error(err))
		} else {
			// 通过父span context开启一个新的span
			span := tracer.StartSpan(info.FullMethod, ext.RPCServerOption(spanContext),
				opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
				ext.SpanKindRPCServer)
			defer span.Finish()
			// 创建新的基于当前span的context
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
		// handler处理rpc请求
		return handler(ctx, req)
	}
}

func NewTracer(serviceName string) (opentracing.Tracer, io.Closer) {
	jConfig := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		// todo add reporter configs
		Reporter: &config.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: "",
			User:              "",
			Password:          "",
		},
	}
	tracer, closer, err := jConfig.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}
