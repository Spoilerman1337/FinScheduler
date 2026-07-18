package domains

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CalendarEvent struct {
	Id          uuid.UUID   `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	Description string      `json:"description" db:"description"`
	Color       string      `json:"color" db:"color"`
	Date        pgtype.Date `json:"date" db:"date"`
}

type CalendarEventTriggerCreate struct {
	Time       pgtype.Time `json:"time"`
	Commentary string      `json:"commentary"`
}

type CalendarEventTriggerUpdate struct {
	Id         *uuid.UUID  `json:"id"`
	Time       pgtype.Time `json:"time"`
	Commentary string      `json:"commentary"`
}

type CalendarEventCreate struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Color       string                       `json:"color"`
	Date        pgtype.Date                  `json:"date"`
	Triggers    []CalendarEventTriggerCreate `json:"triggers"`
}

type CalendarEventUpdate struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Color       string                       `json:"color"`
	Date        pgtype.Date                  `json:"date"`
	Triggers    []CalendarEventTriggerUpdate `json:"triggers"`
}

func (calendarEvent *CalendarEventCreate) Validate() error {
	if !calendarEvent.Date.Valid {
		return fmt.Errorf("date should be valid")
	}

	for _, trigger := range calendarEvent.Triggers {
		if !trigger.Time.Valid {
			return fmt.Errorf("trigger time should be valid")
		}
	}

	return nil
}

func (calendarEvent *CalendarEventUpdate) Validate() error {
	if !calendarEvent.Date.Valid {
		return fmt.Errorf("date should be valid")
	}

	seen := make(map[uuid.UUID]struct{}, len(calendarEvent.Triggers))
	for _, trigger := range calendarEvent.Triggers {
		if !trigger.Time.Valid {
			return fmt.Errorf("trigger time should be valid")
		}

		if trigger.Id == nil {
			continue
		}

		if *trigger.Id == uuid.Nil {
			return fmt.Errorf("trigger id should not be nil")
		}

		if _, exists := seen[*trigger.Id]; exists {
			return fmt.Errorf("trigger id is duplicated: %s", trigger.Id.String())
		}

		seen[*trigger.Id] = struct{}{}
	}

	return nil
}
