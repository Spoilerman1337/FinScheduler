package metrics

import (
	"go.opentelemetry.io/otel/metric"
)

type DatabaseOperation string

const (
	DatabaseOperationNone   DatabaseOperation = "none"
	DatabaseOperationSelect DatabaseOperation = "select"
	DatabaseOperationCount  DatabaseOperation = "count"
	DatabaseOperationInsert DatabaseOperation = "insert"
	DatabaseOperationUpdate DatabaseOperation = "update"
	DatabaseOperationDelete DatabaseOperation = "delete"
)

type Instruments struct {
	HTTPMetrics     HTTPMetrics
	ServiceMetrics  ServiceMetrics
	DatabaseMetrics DatabaseMetrics
}

type HTTPMetrics struct {
	Requests metric.Int64Counter
	Duration metric.Float64Histogram
}

type ServiceMetrics struct {
	Failed metric.Int64Counter
}

type DatabaseMetrics struct {
	Requests metric.Int64Counter
	Duration metric.Float64Histogram
}

var Metrics Instruments

func InitInstruments() {
	Metrics.HTTPMetrics = *newHTTPMetrics(meter)
	Metrics.ServiceMetrics = *newServiceMetrics(meter)
	Metrics.DatabaseMetrics = *newDatabaseMetrics(meter)
}

func newHTTPMetrics(meter metric.Meter) *HTTPMetrics {
	requests, _ := meter.Int64Counter("http_requests_total")
	duration, _ := meter.Float64Histogram("http_request_duration_seconds")

	return &HTTPMetrics{Requests: requests, Duration: duration}
}

func newServiceMetrics(meter metric.Meter) *ServiceMetrics {
	failed, _ := meter.Int64Counter("items_failed_total")

	return &ServiceMetrics{Failed: failed}
}

func newDatabaseMetrics(meter metric.Meter) *DatabaseMetrics {

	requests, _ := meter.Int64Counter("db_requests_total")
	duration, _ := meter.Float64Histogram("db_request_duration_seconds")

	return &DatabaseMetrics{Requests: requests, Duration: duration}
}
