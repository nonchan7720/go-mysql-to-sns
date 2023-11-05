package trace

import (
	"context"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracer struct {
	trace.Tracer
}

func Tracer(t trace.Tracer) trace.Tracer {
	return &tracer{Tracer: t}
}

func (t *tracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	funcName := getFuncName(1)
	if spanName == "" {
		spanName = funcName
	}
	ctx, span := t.Tracer.Start(ctx, spanName, opts...)
	span.SetAttributes(attribute.String("method", funcName))
	return ctx, span
}

func getFuncName(skip int) string {
	pt, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return "unknown_func"
	}
	funcName := runtime.FuncForPC(pt).Name()
	parts := strings.Split(funcName, ".")
	return parts[len(parts)-1]
}
