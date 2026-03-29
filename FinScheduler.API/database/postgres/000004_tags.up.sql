CREATE TABLE tags
(
    id        UUID PRIMARY KEY,
    name      TEXT    NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE tag_to_item
(
    tag_id  UUID REFERENCES tags (id) ON DELETE CASCADE,
    item_id UUID REFERENCES items (id) ON DELETE CASCADE,

    PRIMARY KEY (item_id, tag_id)
);