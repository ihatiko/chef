package tracer

import (
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	commandKey = "command"
)

type Options func(*Tracer)

type Tracer struct {
	ServiceName string
	Command     string
}

func WithCommand(name string) Options {
	return func(tracer *Tracer) {
		tracer.Command = name
	}
}

func WithServiceName(name string) Options {
	return func(tracer *Tracer) {
		tracer.ServiceName = name
	}
}

func (cfg *Config) Use(opts ...Options) {
	tracer := new(Tracer)
	if cfg.Ratio == 0 {
		cfg.Ratio = 0.01
	}

	for _, opt := range opts {
		opt(tracer)
	}
	if tracer.ServiceName == "" {
		tracer.ServiceName = os.Getenv("TECH_SERVICE_NAME")
	}
	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(cfg.Host),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.
					SchemaURL,
				semconv.
					ServiceNameKey.String(tracer.ServiceName),
				attribute.String(
					commandKey, tracer.Command),
			),
		),
		sdktrace.WithSampler(
			sdktrace.TraceIDRatioBased(cfg.Ratio),
		),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{}),
	)
}
