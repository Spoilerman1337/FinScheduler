CREATE TABLE price_forecast_run
(
    id            UUID PRIMARY KEY,
    calculated_at TIMESTAMP NOT NULL,
    item_id       UUID      NOT NULL REFERENCES items (id) ON DELETE CASCADE
);

CREATE INDEX idx_price_forecast_run_item_id
    ON price_forecast_run (item_id);

CREATE TABLE price_forecast
(
    id     UUID PRIMARY KEY,
    run_id UUID           NOT NULL REFERENCES price_forecast_run (id) ON DELETE CASCADE,
    value  NUMERIC(16, 2) NOT NULL CHECK (value >= 0)
);

CREATE INDEX idx_price_forecast_run_id
    ON price_forecast (run_id);
