package traces

import (
	"finscheduler/internal/metrics"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

func RecordHttpSpan(span trace.Span, r *http.Request, route string) {
	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.route", route),
		attribute.String("client.address", r.RemoteAddr),
	)
}

func RecordServiceSpan(span trace.Span, method string) {
	span.SetAttributes(
		attribute.String("service.operation", method),
	)
}

func RecordRepositorySpan(span trace.Span, driver string, operation metrics.DatabaseOperation) {
	span.SetAttributes(
		attribute.String("db.system", driver),
		semconv.DBOperation(string(operation)),
	)
}

func EnrichFailedHttpSpan(span trace.Span, err error, statusCode int) {
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
	span.SetAttributes(attribute.Int("http.status_code", statusCode))
}

func EnrichFailedServiceSpan(span trace.Span, err error) {
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
}

func EnrichFailedRepositorySpanRead(span trace.Span, err error, retrieved int64) {
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
	span.SetAttributes(attribute.Int64("database.rows_retrieved", retrieved))
}

func EnrichFailedRepositorySpanWrite(span trace.Span, err error, affected int64) {
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
	span.SetAttributes(attribute.Int64("database.rows_affected", affected))
}

func EnrichSuccessHttpSpan(span trace.Span, statusCode int) {
	span.SetAttributes(attribute.Int("http.status_code", statusCode))
	span.SetStatus(codes.Ok, "")
}

func EnrichSuccessServiceSpan(span trace.Span) {
	span.SetStatus(codes.Ok, "")
}

func EnrichSuccessRepositorySpanRead(span trace.Span, retrieved int64) {
	span.SetAttributes(attribute.Int64("database.rows_retrieved", retrieved))
	span.SetStatus(codes.Ok, "")
}

func EnrichSuccessRepositorySpanWrite(span trace.Span, affected int64) {
	span.SetAttributes(attribute.Int("http.status_code", statusCode))
}
