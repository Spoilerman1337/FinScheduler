package metrics

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"net/http"
	"time"
)

func RecordHTTPDuration(ctx context.Context, start time.Time) {
	Metrics.HTTPMetrics.Duration.Record(ctx, time.Since(start).Seconds())
}

func RecordHTTPRequest(ctx context.Context, r *http.Request, route string, statusCode int) {
	Metrics.HTTPMetrics.Requests.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("method", r.Method),
			attribute.String("path", r.URL.Path),
			attribute.String("route", route),
			attribute.Int("status", statusCode),
		),
	)
}

func RecordDatabaseDuration(ctx context.Context, start time.Time, driver string, table string, success bool, operation DatabaseOperation) {
	Metrics.DatabaseMetrics.Duration.Record(ctx, time.Since(start).Seconds(),
		metric.WithAttributes(
			attribute.String("driver", driver),
			attribute.String("table", table),
			attribute.Bool("success", success),
			attribute.String("operation", string(operation)),
		),
	)
}

func RecordDatabaseRequest(ctx context.Context, driver string, table string, success bool, operation DatabaseOperation) {
	attributes := []attribute.KeyValue{
		attribute.String("driver", driver),
		attribute.String("table", table),
		attribute.Bool("success", success),
	}
	if operation != DatabaseOperationNone {
		attributes = append(attributes, attribute.String("operation", string(operation)))
	}

	Metrics.DatabaseMetrics.Requests.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("driver", driver),
			attribute.String("table", table),
			attribute.Bool("success", true),
		),
	)
}
