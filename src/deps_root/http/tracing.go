package roothttp

import (
	"context"
	"log"
	"nosebook/src/errors"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func initTracer() (*trace.TracerProvider, error) {
	ctx := context.TODO()
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL("http://jaeger:4318"))
	if err != nil {
		return nil, err
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
	tp, err := errors.Using(initTracer())
	if err != nil {
		log.Fatalln(err)
	}

	this.traceProvider = tp
	this.tracer = tp.Tracer("application")

	this.router.Use(otelgin.Middleware("middleware"))
}

func (this *RootHTTP) HandleCallback(presenterName string) func(name string, ctx context.Context) func() {
	return func(name string, ctx context.Context) func() {
		_, span := this.tracer.Start(ctx, presenterName+"."+name)

		return func() {
			span.End()
		}
	}
}
