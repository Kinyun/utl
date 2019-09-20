package jaeger

import (
	"fmt"
	"github.com/apex/log"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	mtx "github.com/uber/jaeger-lib/metrics"
	prometheus_metrics "github.com/uber/jaeger-lib/metrics/prometheus"
	"io"
	"os"
)

type Optional struct {
	jaegerAdd string
	Name string
}

func Init(logger *log.Entry, opts *Optional) (opentracing.Tracer, io.Closer, error) {
	// Jaeger tracing
	if len(opts.jaegerAdd) < 1{
		opts.jaegerAdd = "127.0.0.1:5775"
	}
	if len(opts.Name) < 1 {
		opts.Name = os.Args[0]
	}

	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 0.05,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: opts.jaegerAdd,
		},
	}
	// Prometheus monitoring
	metrics := prometheus_metrics.New()

	return cfg.New(
		opts.Name,
		config.Logger(jaegerLoggerAdapter{logger}),
		config.Observer(rpcmetrics.NewObserver(metrics.Namespace(mtx.NSOptions{Name: opts.Name, Tags: nil}), rpcmetrics.DefaultNameNormalizer)),
	)
}

type jaegerLoggerAdapter struct {
	logger *log.Entry
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}
