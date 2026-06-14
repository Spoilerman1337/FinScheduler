CREATE TABLE price_history
(
    id          UUID PRIMARY KEY,
    item_id     UUID           NOT NULL REFERENCES items (id) ON DELETE CASCADE,
    recorded_at DATE           NOT NULL,
    value       NUMERIC(16, 2) NOT NULL CHECK (value >= 0)
);

WITH seed_data AS
(
    SELECT
        id AS item_id,
        COALESCE(updated_at, created_at) AS source_recorded_at,
        price AS source_value
    FROM items
)
INSERT INTO price_history (id, item_id, recorded_at, value)
SELECT
    uuidv7(),
    item_id,
    source_recorded_at::date,
    source_value
FROM seed_data;
