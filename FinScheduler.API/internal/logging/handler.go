package logging

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type CustomLoggingHandler struct {
	next slog.Handler
}

func NewCustomLoggingHandler(next slog.Handler) *CustomLoggingHandler {
	return &CustomLoggingHandler{next: next}
}

func (handler *CustomLoggingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return handler.next.Enabled(ctx, level)
}

func (handler *CustomLoggingHandler) Handle(ctx context.Context, record slog.Record) error {
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		record.AddAttrs(
			slog.String("trace_id", spanContext.TraceID().String()),
			slog.String("span_id", spanContext.SpanID().String()),
		)
	}

	return handler.next.Handle(ctx, record)
}

func (handler *CustomLoggingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomLoggingHandler{next: handler.next.WithAttrs(attrs)}
}

func (handler *CustomLoggingHandler) WithGroup(name string) slog.Handler {
	return &CustomLoggingHandler{next: handler.next.WithGroup(name)}
}
