package roothttp

import (
	"context"
	"log"
	"nosebook/src/errors"
	"nosebook/src/lib/config"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func initTracer() (*trace.TracerProvider, error) {
	ctx := context.TODO()
	var exporter trace.SpanExporter
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	if config.Tracing.IsJaegerExporter() {
		exporter, err = otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL("http://jaeger:4318"))
		if err != nil {
			return nil, err
		}
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("nosebook"),
		),
	)

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(r),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func (this *RootHTTP) enableTracing() {
	if !config.Tracing.Enabled() {
		this.tracer = &noop.Tracer{}
		return
	}

	tp, err := errors.Using(initTracer())
	if err != nil {
		log.Fatalln(err)
	}

	this.traceProvider = tp
	this.tracer = tp.Tracer("application")

	this.router.Use(otelgin.Middleware("middleware"))
	this.router.Use(func(ctx *gin.Context) {
		spanCtx := oteltrace.SpanContextFromContext(ctx.Request.Context())
		if spanCtx.HasTraceID() {
			traceId := spanCtx.TraceID()
			ctx.Header("X-Trace-Id", traceId.String())
		}
	})
}
