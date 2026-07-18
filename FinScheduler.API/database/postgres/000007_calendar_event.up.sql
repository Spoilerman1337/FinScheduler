CREATE TABLE calendar_event
(
    id          UUID PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NULL,
    color       TEXT NOT NULL,
    date        DATE NOT NULL
);

CREATE TABLE event_trigger_time
(
    id                UUID PRIMARY KEY,
    time              TIME NOT NULL,
    commentary        TEXT NULL,
    calendar_event_id UUID NOT NULL REFERENCES calendar_event (id) ON DELETE CASCADE,
    CONSTRAINT uq_event_trigger_time_calendar_event_id_time
        UNIQUE (calendar_event_id, time)
);

CREATE INDEX idx_calendar_event_date
    ON calendar_event (date);

CREATE INDEX idx_event_trigger_time_calendar_event_id
    ON event_trigger_time (calendar_event_id);
