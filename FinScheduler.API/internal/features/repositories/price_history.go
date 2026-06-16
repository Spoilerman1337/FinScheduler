package repositories

import (
	"context"
	"finscheduler/internal/features/domains"
	"finscheduler/internal/metrics"
	"finscheduler/internal/traces"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
)

type PriceHistoriesRepository struct {
	db     DBTX
	logger *slog.Logger
}

func NewPriceHistoriesRepository(db DBTX, logger *slog.Logger) *PriceHistoriesRepository {
	return &PriceHistoriesRepository{db: db, logger: logger}
}

func (repository *PriceHistoriesRepository) GetByItemID(ctx context.Context, itemID uuid.UUID) ([]domains.PriceHistory, error) {
	tracer := otel.Tracer("price-histories")
	ctx, span := tracer.Start(ctx, "price-histories-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationSelect)
	defer span.End()

	var priceHistories []domains.PriceHistory

	if itemID == uuid.Nil {
		repository.logger.ErrorContext(ctx, "itemID should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, priceHistoryTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("itemID should not be nil")
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	}

	query := `SELECT recorded_at, value
			  FROM public.price_history
			  WHERE item_id = ?
			  ORDER BY recorded_at DESC`
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(ctx, "executing operation:", "query", query, "itemID", itemID)
	start := time.Now()
	err := sqlx.SelectContext(ctx, repository.db, &priceHistories, query, itemID)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, priceHistoryTableName, err == nil, metrics.DatabaseOperationSelect)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on SELECT operation", "error", err, "itemID", itemID)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, priceHistoryTableName, false, metrics.DatabaseOperationSelect)
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	}

	metrics.RecordDatabaseRequest(ctx, databaseDriver, priceHistoryTableName, true, metrics.DatabaseOperationSelect)
	traces.EnrichSuccessRepositorySpanRead(span, int64(len(priceHistories)))
	return priceHistories, nil
}

func (repository *PriceHistoriesRepository) UpsertToday(ctx context.Context, itemID uuid.UUID, upsert *domains.PriceHistoryUpsert) (*domains.PriceHistory, error) {
	tracer := otel.Tracer("price-histories")
	ctx, span := tracer.Start(ctx, "price-histories-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationUpdate)
	defer span.End()

	if itemID == uuid.Nil {
		repository.logger.ErrorContext(ctx, "itemID should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, priceHistoryTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("itemID should not be nil")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return nil, err
	}

	if upsert == nil {
		repository.logger.ErrorContext(ctx, "upsert should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, priceHistoryTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("upsert should not be nil")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return nil, err
	}

	newID, err := uuid.NewV7()
	if err != nil {
		repository.logger.ErrorContext(ctx, "uuid generation error", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, priceHistoryTableName, false, metrics.DatabaseOperationNone)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return nil, err
	}

	recordedAt := newUTCDate(time.Now().UTC())

	query := `INSERT INTO public.price_history (id, item_id, recorded_at, value)
			  VALUES (?, ?, ?, ?)
			  ON CONFLICT ON CONSTRAINT uq_price_history_item_id_recorded_at
			  DO UPDATE SET value = EXCLUDED.value
			  RETURNING id, item_id, recorded_at, value`
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(ctx, "executing operation:", "query", query, "itemID", itemID, "recordedAt", recordedAt, "value", upsert.Value)
	start := time.Now()
	var priceHistory domains.PriceHistory
	err = sqlx.GetContext(ctx, repository.db, &priceHistory, query, newID, itemID, recordedAt, upsert.Value)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, priceHistoryTableName, err == nil, metrics.DatabaseOperationUpdate)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on UPSERT operation", "error", err, "itemID", itemID, "recordedAt", recordedAt, "value", upsert.Value)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, priceHistoryTableName, false, metrics.DatabaseOperationUpdate)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return nil, err
	}

	metrics.RecordDatabaseRequest(ctx, databaseDriver, priceHistoryTableName, true, metrics.DatabaseOperationUpdate)
	traces.EnrichSuccessRepositorySpanWrite(span, 1)
	return &priceHistory, nil
}

func newUTCDate(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}
