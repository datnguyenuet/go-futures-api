package jaeger

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
	jaegerLog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"io"
)

// InitJaeger Init Jaeger
func InitJaeger() (opentracing.Tracer, io.Closer, error) {
	jaegerCfgInstance := jaegerCfg.Configuration{
		ServiceName: "api_gateway_grpc",
		Sampler: &jaegerCfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerCfg.ReporterConfig{
			LogSpans:           false,
			LocalAgentHostPort: "localhost:6831",
		},
	}

	return jaegerCfgInstance.NewTracer(
		jaegerCfg.Logger(jaegerLog.StdLogger),
		jaegerCfg.Metrics(metrics.NullFactory),
	)
}
