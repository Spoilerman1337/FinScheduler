package domains

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type EventTriggerTime struct {
	Id              uuid.UUID   `json:"id" db:"id"`
	Time            pgtype.Time `json:"time" db:"time"`
	Commentary      string      `json:"commentary" db:"commentary"`
	CalendarEventId uuid.UUID   `json:"calendarEventId" db:"calendar_event_id"`
}

type EventTriggerTimeCreate struct {
	Time            pgtype.Time `json:"time"`
	Commentary      string      `json:"commentary"`
	CalendarEventId uuid.UUID   `json:"calendarEventId"`
}

type EventTriggerTimeUpdate struct {
	Time            pgtype.Time `json:"time"`
	Commentary      string      `json:"commentary"`
	CalendarEventId uuid.UUID   `json:"calendarEventId"`
}
