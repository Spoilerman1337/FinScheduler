CREATE INDEX IF NOT EXISTS idx_items_created_at_id
    ON items (created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_tags_lookup_lower_name_id
    ON tags (LOWER(name), id)
    WHERE is_active = true;
