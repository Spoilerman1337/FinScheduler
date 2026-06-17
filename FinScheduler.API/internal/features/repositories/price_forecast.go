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

type PriceForecastsRepository struct {
	db     DBTX
	logger *slog.Logger
}

func NewPriceForecastsRepository(db DBTX, logger *slog.Logger) *PriceForecastsRepository {
	return &PriceForecastsRepository{db: db, logger: logger}
}

func (repository *PriceForecastsRepository) GetLatestByItemID(ctx context.Context, itemID uuid.UUID) (*domains.PriceForecast, error) {
	tracer := otel.Tracer("price-forecasts")
	ctx, span := tracer.Start(ctx, "price-forecasts-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationSelect)
	defer span.End()

	if itemID == uuid.Nil {
		repository.logger.ErrorContext(ctx, "itemID should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, priceForecastTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("itemID should not be nil")
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	}

	query := `SELECT id, item_id, calculated_at, last_known_price, average_monthly_drift
			  FROM public.price_forecast
			  WHERE item_id = ?
			  ORDER BY calculated_at DESC, id DESC
			  LIMIT 1`
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(ctx, "executing operation:", "query", query, "itemID", itemID)
	start := time.Now()
	var priceForecast domains.PriceForecast
	err := sqlx.GetContext(ctx, repository.db, &priceForecast, query, itemID)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, priceForecastTableName, err == nil, metrics.DatabaseOperationSelect)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on SELECT operation", "error", err, "itemID", itemID)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, priceForecastTableName, false, metrics.DatabaseOperationSelect)
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	}

	metrics.RecordDatabaseRequest(ctx, databaseDriver, priceForecastTableName, true, metrics.DatabaseOperationSelect)
	traces.EnrichSuccessRepositorySpanRead(span, 1)
	return &priceForecast, nil
}

func (repository *PriceForecastsRepository) UpsertByItemID(ctx context.Context, itemID uuid.UUID, upsert *domains.PriceForecastUpsert) (*domains.PriceForecast, error) {
	tracer := otel.Tracer("price-forecasts")
	ctx, span := tracer.Start(ctx, "price-forecasts-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationUpdate)
	defer span.End()

	if itemID == uuid.Nil {
		repository.logger.ErrorContext(ctx, "itemID should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, priceForecastTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("itemID should not be nil")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return nil, err
	}

	if upsert == nil {
		repository.logger.ErrorContext(ctx, "upsert should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, priceForecastTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("upsert should not be nil")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return nil, err
	}

	newID, err := uuid.NewV7()
	if err != nil {
		repository.logger.ErrorContext(ctx, "uuid generation error", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, priceForecastTableName, false, metrics.DatabaseOperationNone)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return nil, err
	}

	calculatedAt := newUTCDate(upsert.CalculatedAt.UTC())

	query := `INSERT INTO public.price_forecast (id, item_id, calculated_at, last_known_price, average_monthly_drift)
			  VALUES (?, ?, ?, ?, ?)
			  ON CONFLICT ON CONSTRAINT uq_price_forecast_item_id_calculated_at
			  DO UPDATE SET
			      last_known_price = EXCLUDED.last_known_price,
			      average_monthly_drift = EXCLUDED.average_monthly_drift
			  RETURNING id, item_id, calculated_at, last_known_price, average_monthly_drift`
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(
		ctx,
		"executing operation:",
		"query", query,
		"itemID", itemID,
		"calculatedAt", calculatedAt,
		"lastKnownPrice", upsert.LastKnownPrice,
		"averageMonthlyDrift", upsert.AverageMonthlyDrift,
	)
	start := time.Now()
	var priceForecast domains.PriceForecast
	err = sqlx.GetContext(
		ctx,
		repository.db,
		&priceForecast,
		query,
		newID,
		itemID,
		calculatedAt,
		upsert.LastKnownPrice,
		upsert.AverageMonthlyDrift,
	)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, priceForecastTableName, err == nil, metrics.DatabaseOperationUpdate)
	if err != nil {
		repository.logger.ErrorContext(
			ctx,
			"error on UPSERT operation",
			"error", err,
			"itemID", itemID,
			"calculatedAt", calculatedAt,
			"lastKnownPrice", upsert.LastKnownPrice,
			"averageMonthlyDrift", upsert.AverageMonthlyDrift,
		)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, priceForecastTableName, false, metrics.DatabaseOperationUpdate)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return nil, err
	}

	metrics.RecordDatabaseRequest(ctx, databaseDriver, priceForecastTableName, true, metrics.DatabaseOperationUpdate)
	traces.EnrichSuccessRepositorySpanWrite(span, 1)
	return &priceForecast, nil
}
