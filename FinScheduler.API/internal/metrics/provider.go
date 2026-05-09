package metrics

import (
	"context"
	"finscheduler/internal/infra"
	"fmt"
	"net/http"

	prometheusclient "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	prometheusexporter "go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var meter = otel.Meter("fin-scheduler-api")
var metricsHandler http.Handler = http.NotFoundHandler()

func InitMetrics(ctx context.Context, cfg *infra.Config) (*sdkmetric.MeterProvider, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.Observability.ServiceName),
			semconv.DeploymentEnvironmentKey.String(cfg.Env),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create metrics resource: %w", err)
	}

	if !cfg.Observability.Metrics.Enabled {
		mp := sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(res),
			sdkmetric.WithReader(sdkmetric.NewManualReader()),
		)
		otel.SetMeterProvider(mp)
		meter = mp.Meter(cfg.Observability.ServiceName)
		metricsHandler = http.NotFoundHandler()

		return mp, nil
	}

	registry := prometheusclient.NewRegistry()
	exporter, err := prometheusexporter.New(
		prometheusexporter.WithRegisterer(registry),
	)
	if err != nil {
		return nil, fmt.Errorf("create prometheus exporter: %w", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(exporter),
	)
	otel.SetMeterProvider(mp)
	meter = mp.Meter(cfg.Observability.ServiceName)
	metricsHandler = promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	return mp, nil
}

func Handler() http.Handler {
	return metricsHandler
}
