package tracking

import (
	"context"
	"log/slog"
	"time"

	"github.com/nonchan7720/go-storage-to-messenger/pkg/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func NewTracerProvider(
	ctx context.Context,
	enabled bool,
	otelAgentAddr string,
	serviceName string,
	environment string,
	_serviceRoot string,
) (trace.TracerProvider, func(), error) {
	var (
		tp      trace.TracerProvider
		cleanup func()
	)
	serviceRoot = _serviceRoot
	if !enabled {
		tp = trace.NewNoopTracerProvider()
		cleanup = func() {}
	} else {
		traceClient := otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(otelAgentAddr),
			otlptracegrpc.WithDialOption(grpc.WithBlock()),
		)
		timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		exporter, err := otlptrace.New(timeoutCtx, traceClient)
		if err != nil {
			return nil, nil, err
		}
		r := NewResource(serviceName, "v1", environment)
		sdkTP := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(r),
		)
		pp := NewPropagator()
		otel.SetTextMapPropagator(pp)
		cleanup = func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := sdkTP.ForceFlush(ctx); err != nil {
				slog.With(logging.WithErr(err)).Error(err.Error())
			}
			ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
			if err := sdkTP.Shutdown(ctx2); err != nil {
				slog.With(logging.WithErr(err)).Error(err.Error())
			}
			cancel()
			cancel2()
		}
		tp = sdkTP
	}
	otel.SetTracerProvider(tp)
	return tp, cleanup, nil
}

func NewResource(serviceName string, version string, environment string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(version),
		semconv.DeploymentEnvironmentKey.String(environment),
		attribute.String("environment", environment),
		attribute.String("env", environment),
	)
}

func NewPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
