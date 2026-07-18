package domains

import "github.com/google/uuid"

type EventTriggerTime struct {
	Id              uuid.UUID `json:"id" db:"id"`
	Time            TimeOnly  `json:"time" db:"time"`
	Commentary      string    `json:"commentary" db:"commentary"`
	CalendarEventId uuid.UUID `json:"calendarEventId" db:"calendar_event_id"`
}

type EventTriggerTimeCreate struct {
	Time            TimeOnly  `json:"time"`
	Commentary      string    `json:"commentary"`
	CalendarEventId uuid.UUID `json:"calendarEventId"`
}

type EventTriggerTimeUpdate struct {
	Time            TimeOnly  `json:"time"`
	Commentary      string    `json:"commentary"`
	CalendarEventId uuid.UUID `json:"calendarEventId"`
}
