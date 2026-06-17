CREATE TABLE price_forecast
(
    id                    UUID PRIMARY KEY,
    item_id               UUID           NOT NULL REFERENCES items (id) ON DELETE CASCADE,
    calculated_at         TIMESTAMP      NOT NULL,
    last_known_price      NUMERIC(16, 2) NOT NULL CHECK (last_known_price >= 0),
    average_monthly_drift NUMERIC(16, 2) NOT NULL
);

CREATE INDEX idx_price_forecast_item_id
    ON price_forecast (item_id);
