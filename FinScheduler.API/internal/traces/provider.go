package traces

import (
	"context"
	"finscheduler/internal/infra"
	"fmt"
	"net/url"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func InitTracer(ctx context.Context, cfg *infra.Config) (*sdktrace.TracerProvider, error) {
	options := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.Observability.ServiceName),
			attribute.String("env", cfg.Env),
		)),
	}

	if cfg.Observability.Traces.Enabled {
		exporter, err := newTraceExporter(ctx, cfg)
		if err != nil {
			return nil, err
		}

		options = append(options,
			sdktrace.WithBatcher(exporter),
			sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.Observability.Traces.RootTraceSamplingRatio))),
		)
	} else {
		options = append(options, sdktrace.WithSampler(sdktrace.NeverSample()))
	}

	tp := sdktrace.NewTracerProvider(options...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

func newTraceExporter(ctx context.Context, cfg *infra.Config) (*otlptrace.Exporter, error) {
	options, err := newTraceExporterOptions(cfg.Observability.Traces.ExportEndpoint)
	if err != nil {
		return nil, err
	}

	client := otlptracehttp.NewClient(options...)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("create trace exporter: %w", err)
	}

	return exporter, nil
}

func newTraceExporterOptions(rawExportEndpoint string) ([]otlptracehttp.Option, error) {
	exportEndpoint := strings.TrimSpace(rawExportEndpoint)
	if exportEndpoint == "" {
		return nil, fmt.Errorf("trace export endpoint is empty")
	}

	parsedEndpoint, err := url.Parse(exportEndpoint)
	if err != nil {
		return nil, fmt.Errorf("parse trace export endpoint: %w", err)
	}
	if parsedEndpoint.Scheme == "" {
		return nil, fmt.Errorf("trace export endpoint must include http:// or https://")
	}
	if parsedEndpoint.Host == "" {
		return nil, fmt.Errorf("trace export endpoint host is empty")
	}
	scheme := strings.ToLower(parsedEndpoint.Scheme)
	if scheme != "http" && scheme != "https" {
		return nil, fmt.Errorf("unsupported trace export endpoint scheme %q", parsedEndpoint.Scheme)
	}

	options := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(parsedEndpoint.Host),
	}
	if scheme == "http" {
		options = append(options, otlptracehttp.WithInsecure())
	}

	if parsedEndpoint.Path != "" && parsedEndpoint.Path != "/" {
		options = append(options, otlptracehttp.WithURLPath(parsedEndpoint.Path))
	}

	return options, nil
}
