package interceptors

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type InterceptorManager struct {
	logger *zap.SugaredLogger
	tracer opentracing.Tracer
}

// NewInterceptorManager InterceptorManager constructor
func NewInterceptorManager(logger *zap.SugaredLogger, tracer opentracing.Tracer) *InterceptorManager {
	return &InterceptorManager{logger: logger, tracer: tracer}
}

// Logger Interceptor
func (im *InterceptorManager) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	im.logger.Infof("Method: %s, Time: %v, Metadata: %v, Err: %v", info.FullMethod, time.Since(start), md, err)

	return reply, err
}

func (im *InterceptorManager) GetInterceptor() func(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		im.logger.Infof(
			"call=%v req=%#v reply=%#v time=%v err=%v", method, req, reply, time.Since(start), err,
		)
		return err
	}
}

func (im *InterceptorManager) GetTracer() opentracing.Tracer {
	return im.tracer
}
